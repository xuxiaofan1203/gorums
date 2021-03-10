// Package gengorums is internal to the gorums protobuf module.
package gengorums

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/relab/gorums"
	"github.com/relab/gorums/internal/correctable"
	"github.com/relab/gorums/internal/version"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// GenerateVersionMarkers specifies whether to generate version markers.
var GenerateVersionMarkers = true

// GenerateFile generates a _gorums.pb.go file containing Gorums service definitions.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) {
	if !gorumsGuard(file) {
		return
	}
	filename := file.GeneratedFilenamePrefix + "_gorums.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	genGeneratedHeader(gen, g, file)
	g.P("package ", file.GoPackageName)
	g.P()
	genVersionCheck(g)
	g.P(staticCode)
	g.P()
	for path, ident := range pkgIdentMap {
		addImport(path, ident, g)
	}
	generateFileContent(gen, file, g)
}

func genGeneratedHeader(gen *protogen.Plugin, g *protogen.GeneratedFile, f *protogen.File) {
	g.P("// Code generated by protoc-gen-gorums. DO NOT EDIT.")

	if GenerateVersionMarkers {
		g.P("// versions:")
		protocVersion := "(unknown)"
		if v := gen.Request.GetCompilerVersion(); v != nil {
			protocVersion = fmt.Sprintf("v%v.%v.%v", v.GetMajor(), v.GetMinor(), v.GetPatch())
		}
		g.P("// \tprotoc-gen-gorums ", version.String())
		g.P("// \tprotoc            ", protocVersion)
	}

	if f.Proto.GetOptions().GetDeprecated() {
		g.P("// ", f.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", f.Desc.Path())
	}
	g.P()
}

func genVersionCheck(g *protogen.GeneratedFile) {
	if GenerateVersionMarkers {
		g.P("const (")
		g.P("// Verify that this generated code is sufficiently up-to-date.")
		g.P("_ = gorums.EnforceVersion(", gorums.GenVersion, " - gorums.MinVersion)")
		g.P("// Verify that the gorums runtime is sufficiently up-to-date.")
		g.P("_ = gorums.EnforceVersion(gorums.MaxVersion - ", gorums.GenVersion, ")")
		g.P(")")
		g.P()
	}
}

// gorumsGuard returns true if there is something for Gorums to generate
// for the given file. If it returns false, there is nothing for Gorums
// generate for this file. If may also fail with an error message.
func gorumsGuard(file *protogen.File) bool {
	if len(file.Services) == 0 || !hasGorumsMethods(file.Services) {
		// there is nothing for this plugin to do
		return false
	}
	if len(file.Services) > 1 {
		// To build multiple services, make separate proto files and
		// run the plugin separately for each proto file.
		// These cannot share the same Go package.
		log.Fatalln("Gorums does not support multiple services in the same proto file.")
	}
	// fail generator if a Gorums reserved identifier is used as a message name.
	for _, msg := range file.Messages {
		msgName := fmt.Sprintf("%v", msg.Desc.Name())
		for _, reserved := range reservedIdents {
			if msgName == reserved {
				log.Fatalf("%v.proto: contains message %s, which is a reserved Gorums type.\n", file.GeneratedFilenamePrefix, msgName)
			}
		}
	}
	return true
}

// GenerateFileContent generates the Gorums service definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	// sort the gorums types so that output remains stable across rebuilds
	sortedTypes := make([]string, 0, len(gorumsCallTypesInfo))
	for gorumsType := range gorumsCallTypesInfo {
		sortedTypes = append(sortedTypes, gorumsType)
	}
	sort.Strings(sortedTypes)
	for _, gorumsType := range sortedTypes {
		genGorumsType(g, file.Services, gorumsType)
	}
}

// servicesData hold the services to generate and a reference to the file in which
// the services should be generated. This is data to be used by template generator.
type servicesData struct {
	GenFile  *protogen.GeneratedFile
	Services []*protogen.Service
}

// genGorumsType generates Gorums methods and corresponding datastructures for gorumsType.
func genGorumsType(g *protogen.GeneratedFile, services []*protogen.Service, gorumsType string) {
	data := servicesData{g, services}
	if callTypeInfo := gorumsCallTypesInfo[gorumsType]; callTypeInfo.extInfo == nil {
		g.P(mustExecute(parseTemplate(gorumsType, callTypeInfo.template), data))
	} else {
		genGorumsMethods(gorumsType, data, callTypeInfo)
	}
	g.P()
}

// genGorumsMethods generates Gorums methods for the given call type.
func genGorumsMethods(gorumsType string, data servicesData, callTypeInfo *callTypeInfo) {
	type methodData struct {
		GenFile *protogen.GeneratedFile
		Method  *protogen.Method
	}
	g := data.GenFile
	for _, service := range data.Services {
		for _, method := range service.Methods {
			err := validateOptions(method)
			if err != nil {
				log.Fatal(err)
			}
			callType := callTypeInfo.deriveCallType(method)
			if callType.check(method) {
				fmt.Fprintf(os.Stderr, "generating(%v): %s\n", gorumsType, method.GoName)
				g.P(mustExecute(parseTemplate(gorumsType, callType.template), methodData{g, method}))
			}
		}
	}
}

// callType returns call type information for the given method.
// If the given method has specified a Gorums method option that
// correspond to a call type, this call type is returned. Further,
// if the call type has a sub call type, then this is returned instead.
func callType(method *protogen.Method) *callTypeInfo {
	for _, callTypeInfo := range callTypeOptions(method) {
		callType := callTypeInfo.deriveCallType(method)
		if callType.check(method) {
			return callType
		}
	}
	panic(fmt.Sprintf("unknown call type for method %s\n", method.GoName))
}

// callTypeOptions returns all Gorums call types for the given method.
func callTypeOptions(method *protogen.Method) []*callTypeInfo {
	methExt := protoimpl.X.MessageOf(method.Desc.Options()).Interface()
	var options []*callTypeInfo
	for _, callType := range gorumsCallTypesInfo {
		if callType.extInfo != nil {
			if proto.HasExtension(methExt, callType.extInfo) {
				options = append(options, callType)
			}
		}
	}
	return options
}

// hasGorumsMethods returns true if one of the methods in the set of services
// has a Gorums call type.
func hasGorumsMethods(services []*protogen.Service) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			for _, callTypeInfo := range gorumsCallTypesInfo {
				callType := callTypeInfo.deriveCallType(method)
				if callType.check(method) {
					return true
				}
			}
		}
	}
	return false
}

// callTypeInfo holds information about the option type, the option type name,
// documentation string for the option type, the template used to generate
// a method annotated with the given option, and a chkFn function that returns
// true if code for the option type should be generated for the given method.
type callTypeInfo struct {
	extInfo        *protoimpl.ExtensionInfo
	docName        string
	template       string
	outPrefix      string
	chkFn          func(m *protogen.Method) bool
	nestedCallType map[string]*callTypeInfo
}

// check returns true if the given method is associated with this call type.
func (c *callTypeInfo) check(m *protogen.Method) bool {
	if c != nil && c.chkFn != nil {
		return c.chkFn(m)
	}
	return false
}

// deriveCallType resolves the nested call type if any.
func (c *callTypeInfo) deriveCallType(m *protogen.Method) *callTypeInfo {
	if c != nil {
		if c.nestedCallType != nil {
			for _, nestedCallType := range c.nestedCallType {
				if nestedCallType.chkFn(m) {
					return nestedCallType
				}
			}
		}
	}
	return c
}

// callTypeName returns the name of the call type option with the prefix removed.
func callTypeName(ext *protoimpl.ExtensionInfo) string {
	s := string(ext.TypeDescriptor().FullName())
	return s[strings.LastIndex(s, ".")+1:]
}

// gorumsCallTypesInfo maps Gorums call type names to callTypeInfo.
// This includes details such as the template, extension info and
// a chkFn function used to check for the particular call type.
// The entries in this map is used to generate dev/zorums_{type}.pb.go
// files for the different keys.
var gorumsCallTypesInfo = map[string]*callTypeInfo{
	"qspec":  {template: qspecInterface},
	"types":  {template: datatypes},
	"server": {template: server},

	callTypeName(gorums.E_Rpc): {
		extInfo:  gorums.E_Rpc,
		docName:  "rpc",
		template: rpcCall,
		chkFn: func(m *protogen.Method) bool {
			return !hasMethodOption(m, gorumsCallTypes...)
		},
	},
	callTypeName(gorums.E_Quorumcall): {
		extInfo:  gorums.E_Quorumcall,
		docName:  "quorum",
		template: quorumCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Quorumcall) && !hasMethodOption(m, gorums.E_Async)
		},
	},
	callTypeName(gorums.E_Async): {
		extInfo:   gorums.E_Async,
		docName:   "asynchronous quorum",
		template:  asyncCall,
		outPrefix: "Async",
		chkFn: func(m *protogen.Method) bool {
			return hasAllMethodOption(m, gorums.E_Quorumcall, gorums.E_Async)
		},
	},
	callTypeName(gorums.E_Correctable): {
		extInfo: gorums.E_Correctable,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Correctable)
		},
		nestedCallType: map[string]*callTypeInfo{
			callTypeName(correctable.E_Correctable): {
				extInfo:   correctable.E_Correctable,
				docName:   "correctable quorum",
				template:  correctableCall,
				outPrefix: "Correctable",
				chkFn: func(m *protogen.Method) bool {
					return hasMethodOption(m, gorums.E_Correctable) && !m.Desc.IsStreamingServer()
				},
			},
			callTypeName(correctable.E_CorrectableStream): {
				extInfo:   correctable.E_CorrectableStream,
				docName:   "correctable stream quorum",
				template:  correctableCall,
				outPrefix: "CorrectableStream",
				chkFn: func(m *protogen.Method) bool {
					return hasMethodOption(m, gorums.E_Correctable) && m.Desc.IsStreamingServer()
				},
			},
		},
	},
	callTypeName(gorums.E_Multicast): {
		extInfo:  gorums.E_Multicast,
		docName:  "multicast",
		template: multicastCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Multicast)
		},
	},
	callTypeName(gorums.E_Unicast): {
		extInfo:  gorums.E_Unicast,
		docName:  "unicast",
		template: unicastCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Unicast)
		},
	},
}

// gorumsCallTypes should list all available call types supported by Gorums.
// These are considered mutually incompatible.
var gorumsCallTypes = []*protoimpl.ExtensionInfo{
	gorums.E_Quorumcall,
	gorums.E_Async,
	gorums.E_Correctable,
	gorums.E_Multicast,
	gorums.E_Unicast,
}

// callTypesWithInternal should list all available call types that
// has a quorum function and hence need an internal type that wraps
// the return type with additional information.
var callTypesWithInternal = []*protoimpl.ExtensionInfo{
	gorums.E_Quorumcall,
	gorums.E_Async,
	gorums.E_Correctable,
}

// callTypesWithPromiseObject lists all call types that returns
// a promise (async or correctable) object.
var callTypesWithPromiseObject = []*protoimpl.ExtensionInfo{
	gorums.E_Async,
	gorums.E_Correctable,
}

// hasGorumsCallType returns true if the given method has specified
// one of the call types supported by Gorums.
func hasGorumsCallType(method *protogen.Method) bool {
	return hasMethodOption(method, gorumsCallTypes...)
}

// hasMethodOption returns true if the method has one of the given method options.
func hasMethodOption(method *protogen.Method, methodOptions ...*protoimpl.ExtensionInfo) bool {
	ext := protoimpl.X.MessageOf(method.Desc.Options()).Interface()
	for _, callType := range methodOptions {
		if proto.HasExtension(ext, callType) {
			return true
		}
	}
	return false
}

// hasAllMethodOption returns true if the method has all of the given method options
func hasAllMethodOption(method *protogen.Method, methodOptions ...*protoimpl.ExtensionInfo) bool {
	ext := protoimpl.X.MessageOf(method.Desc.Options()).Interface()
	for _, callType := range methodOptions {
		if !proto.HasExtension(ext, callType) {
			return false
		}
	}
	return true
}

// validateOptions returns an error if the extension options
// for the provided method are incompatible.
func validateOptions(method *protogen.Method) error {
	switch {
	case hasMethodOption(method, gorums.E_Async) && !hasMethodOption(method, gorums.E_Quorumcall):
		return optionErrorf("is required for async methods", method, gorums.E_Quorumcall)

	case !hasMethodOption(method, gorums.E_Multicast) && method.Desc.IsStreamingClient():
		return optionErrorf("is required for client-server stream methods", method, gorums.E_Multicast)

	case !hasMethodOption(method, gorums.E_Correctable) && method.Desc.IsStreamingServer():
		return optionErrorf("is required for server-client stream methods", method, gorums.E_Correctable)

	case hasMethodOption(method, gorums.E_Correctable) && method.Desc.IsStreamingClient():
		return optionErrorf("is only valid for server-client stream methods", method, gorums.E_Correctable)
	}
	return nil
}

func optionErrorf(s string, method *protogen.Method, ext *protoimpl.ExtensionInfo) error {
	return fmt.Errorf("%s.%s: option '%s' "+s, method.Parent.Desc.Name(), method.Desc.Name(), ext.TypeDescriptor().FullName())
}
