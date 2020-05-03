// Package gengorums is internal to the gorums protobuf module.
package gengorums

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/relab/gorums"
	"github.com/relab/gorums/internal/ordering"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// TODO(meling) replace github.com/relab/gorums with gorums.io as import package

// GenerateFile generates a _gorums.pb.go file containing Gorums service definitions.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 || !checkMethods(file.Services, func(m *protogen.Method) bool { return hasGorumsCallType(m) }) {
		// there is nothing for this plugin to do
		return nil
	}
	if len(file.Services) > 1 {
		// To build multiple services, make separate proto files and
		// run the plugin separately for each proto file.
		// These cannot share the same Go package.
		log.Fatalln("Gorums does not support multiple services in the same proto file.")
	}
	// TODO(meling) make this more generic; figure out what are the reserved types from the static files.
	for _, msg := range file.Messages {
		msgName := fmt.Sprintf("%v", msg.Desc.Name())
		for _, reserved := range []string{"Configuration", "Node", "Manager", "ManagerOption"} {
			if msgName == reserved {
				log.Fatalf("%v.proto: contains message %s, which is a reserved Gorums type.\n", file.GeneratedFilenamePrefix, msgName)
			}
		}
	}

	filename := file.GeneratedFilenamePrefix + "_gorums.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-gorums. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P(staticCode)
	g.P()
	for path, ident := range pkgIdentMap {
		addImport(path, ident, g)
	}
	GenerateFileContent(gen, file, g)
	return g
}

// GenerateFileContent generates the Gorums service definitions, excluding the package statement.
func GenerateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	data := servicesData{g, file.Services}
	// sort the gorums types so that output remains stable across rebuilds
	sortedTypes := make([]string, 0, len(gorumsCallTypesInfo))
	for gorumsType := range gorumsCallTypesInfo {
		sortedTypes = append(sortedTypes, gorumsType)
	}
	sort.Strings(sortedTypes)
	for _, gorumsType := range sortedTypes {
		if callTypeInfo := gorumsCallTypesInfo[gorumsType]; callTypeInfo.extInfo == nil {
			g.P(mustExecute(parseTemplate(gorumsType, callTypeInfo.template), data))
		} else {
			genGorumsMethods(data, callTypeInfo.extInfo)
		}
		g.P()
	}
}

func genGorumsMethods(data servicesData, methodOptions ...*protoimpl.ExtensionInfo) {
	g := data.GenFile
	for _, service := range data.Services {
		for _, method := range service.Methods {
			if hasMethodOption(method, gorums.E_Ordered) {
				if hasOrderingOption(method, methodOptions...) {
					fmt.Fprintf(os.Stderr, "processing %s\n", method.GoName)
					g.P(genGorumsMethod(g, method))
				}
			} else if hasMethodOption(method, methodOptions...) {
				fmt.Fprintf(os.Stderr, "processing %s\n", method.GoName)
				g.P(genGorumsMethod(g, method))
			}
		}
	}
}

func genGorumsMethod(g *protogen.GeneratedFile, method *protogen.Method) string {
	callTypeInfo := methodCallTypeInfo(method)
	return mustExecute(parseTemplate(callTypeInfo.optionName, callTypeInfo.template), methodData{g, method})
}

func callTypeName(method *protogen.Method) string {
	return methodCallTypeInfo(method).docName
}

func methodCallTypeInfo(method *protogen.Method) *callTypeInfo {
	methodOption := validateMethodExtensions(method)
	if optionName, ok := reverseMap[methodOption]; ok {
		if callTypeInfo, ok := gorumsCallTypesInfo[optionName]; ok {
			return callTypeInfo
		}
	}
	panic(fmt.Sprintf("unknown method type %s\n", method.GoName))
}

type servicesData struct {
	GenFile  *protogen.GeneratedFile
	Services []*protogen.Service
}

type methodData struct {
	GenFile *protogen.GeneratedFile
	Method  *protogen.Method
}

// hasGorumsType returns true if one of the service methods specify
// the given gorums type.
func hasGorumsType(services []*protogen.Service, gorumsType string) bool {
	if callTypeInfo := gorumsCallTypesInfo[gorumsType]; callTypeInfo.extInfo == nil {
		return true
	}
	return checkMethodOptions(services, gorumsType)
}

// checkMethodOptions returns true if one of the service methods defines
// the given gorums option type.
func checkMethodOptions(services []*protogen.Service, option string) bool {
	if callTypeInfo, ok := gorumsCallTypesInfo[option]; ok {
		return checkMethods(services, func(m *protogen.Method) bool {
			return callTypeInfo.chkFn(m)
		})
	}
	return false
}

// checkMethods returns true if the function fn evaluates to true
// for one of the methods in the set of services.
func checkMethods(services []*protogen.Service, fn func(m *protogen.Method) bool) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			if fn(method) {
				return true
			}
		}
	}
	return false
}

// compute index to start of option name
const index = len("gorums.")
const soIndex = len("ordering.")

// mapping from option type to option name
var reverseMap = map[*protoimpl.ExtensionInfo]string{
	gorums.E_Quorumcall:        gorums.E_Quorumcall.Name[index:],
	gorums.E_QcFuture:          gorums.E_QcFuture.Name[index:],
	gorums.E_Correctable:       gorums.E_Correctable.Name[index:],
	gorums.E_CorrectableStream: gorums.E_CorrectableStream.Name[index:],
	gorums.E_Multicast:         gorums.E_Multicast.Name[index:],
	ordering.E_OrderedQc:       ordering.E_OrderedQc.Name[soIndex:],
	ordering.E_OrderedRpc:      ordering.E_OrderedRpc.Name[soIndex:],
}

// callTypeInfo holds information about the option type, the option type name,
// documentation string for the option type, the template used to generate
// a method annotated with the given option, and a chkFn function that returns
// true if code for the option type should be generated for the given method.
type callTypeInfo struct {
	extInfo    *protoimpl.ExtensionInfo
	optionName string
	docName    string
	template   string
	chkFn      func(m *protogen.Method) bool
}

// gorumsCallTypesInfo maps Gorums call type names to callTypeInfo.
// This includes details such as the template, extension info and
// a chkFn function used to check for the particular call type.
// The entries in this map is used to generate dev/zorums_{type}.pb.go
// files for the different keys.
var gorumsCallTypesInfo = map[string]*callTypeInfo{
	"node":  {template: node},
	"qspec": {template: qspecInterface},
	"types": {template: datatypes},

	gorums.E_Quorumcall.Name[index:]: {
		extInfo:    gorums.E_Quorumcall,
		optionName: gorums.E_Quorumcall.Name[index:],
		docName:    "quorum",
		template:   quorumCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Quorumcall)
		},
	},
	gorums.E_QcFuture.Name[index:]: {
		extInfo:    gorums.E_QcFuture,
		optionName: gorums.E_QcFuture.Name[index:],
		docName:    "asynchronous quorum",
		template:   futureCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_QcFuture)
		},
	},
	gorums.E_Correctable.Name[index:]: {
		extInfo:    gorums.E_Correctable,
		optionName: gorums.E_Correctable.Name[index:],
		docName:    "correctable quorum",
		template:   correctableCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Correctable)
		},
	},
	gorums.E_CorrectableStream.Name[index:]: {
		extInfo:    gorums.E_CorrectableStream,
		optionName: gorums.E_CorrectableStream.Name[index:],
		docName:    "correctable stream quorum",
		template:   correctableStreamCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_CorrectableStream)
		},
	},
	gorums.E_Multicast.Name[index:]: {
		extInfo:    gorums.E_Multicast,
		optionName: gorums.E_Multicast.Name[index:],
		docName:    "multicast",
		template:   multicastCall,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Multicast)
		},
	},
	ordering.E_OrderedQc.Name[soIndex:]: {
		extInfo:    ordering.E_OrderedQc,
		optionName: ordering.E_OrderedQc.Name[soIndex:],
		docName:    "ordered quorum",
		template:   orderingQC,
		chkFn: func(m *protogen.Method) bool {
			return hasAllMethodOption(m, gorums.E_Ordered, gorums.E_Quorumcall)
		},
	},
	ordering.E_OrderedRpc.Name[soIndex:]: {
		extInfo:    ordering.E_OrderedRpc,
		optionName: ordering.E_OrderedRpc.Name[soIndex:],
		docName:    "ordered",
		template:   orderingRPC,
		chkFn: func(m *protogen.Method) bool {
			return hasMethodOption(m, gorums.E_Ordered) && !hasGorumsCallType(m)
		},
	},
}

// mapping from ordering type to a checker that will check if a method has that type
var orderingTypeCheckers = map[*protoimpl.ExtensionInfo]func(*protogen.Method) bool{
	ordering.E_OrderedQc: func(m *protogen.Method) bool {
		return hasAllMethodOption(m, gorums.E_Ordered, gorums.E_Quorumcall)
	},
	ordering.E_OrderedRpc: func(m *protogen.Method) bool {
		return hasMethodOption(m, gorums.E_Ordered) && !hasGorumsCallType(m)
	},
}

// gorumsCallTypes should list all available call types supported by Gorums.
// These are considered mutually incompatible.
var gorumsCallTypes = []*protoimpl.ExtensionInfo{
	gorums.E_Quorumcall,
	gorums.E_QcFuture,
	gorums.E_Correctable,
	gorums.E_CorrectableStream,
	gorums.E_Multicast,
}

// callTypesWithInternal should list all available call types that
// has a quorum function and hence need an internal type that wraps
// the return type with additional information.
var callTypesWithInternal = []*protoimpl.ExtensionInfo{
	gorums.E_Quorumcall,
	gorums.E_QcFuture,
	gorums.E_Correctable,
	gorums.E_CorrectableStream,
}

// callTypesWithPromiseObject lists all call types that returns
// a promise (future or correctable) object.
var callTypesWithPromiseObject = []*protoimpl.ExtensionInfo{
	gorums.E_QcFuture,
	gorums.E_Correctable,
	gorums.E_CorrectableStream,
}

// orderingCallTypes should list call types that use ordering.
var orderingCallTypes = []*protoimpl.ExtensionInfo{
	ordering.E_OrderedQc,
	ordering.E_OrderedRpc,
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

// hasOrderingOption returns true if the method has one of the ordering method options.
func hasOrderingOption(method *protogen.Method, methodOptions ...*protoimpl.ExtensionInfo) bool {
	for _, option := range methodOptions {
		if f, ok := orderingTypeCheckers[option]; ok && f(method) {
			return true
		}
	}
	return false
}

// validateMethodExtensions returns the method option for the
// call type of the given method. If the method specifies multiple
// call types, validation will fail with a panic.
func validateMethodExtensions(method *protogen.Method) *protoimpl.ExtensionInfo {
	methExt := protoimpl.X.MessageOf(method.Desc.Options()).Interface()
	var firstOption *protoimpl.ExtensionInfo
	for _, callType := range gorumsCallTypes {
		if proto.HasExtension(methExt, callType) {
			if firstOption != nil {
				log.Fatalf("%s.%s: cannot combine options: '%s' and '%s'",
					method.Parent.Desc.Name(), method.Desc.Name(), firstOption.Name, callType.Name)
			}
			firstOption = callType
		}
	}

	// check if the method matches any ordering types
	for t, f := range orderingTypeCheckers {
		if f(method) {
			firstOption = t
		}
	}

	isQuorumCallVariant := hasMethodOption(method, callTypesWithInternal...)
	switch {
	case !isQuorumCallVariant && proto.GetExtension(methExt, gorums.E_CustomReturnType) != "":
		// Only QC variants can define custom return type
		// (we don't support rewriting the plain gRPC methods.)
		log.Fatalf(
			"%s.%s: cannot combine non-quorum call method with the '%s' option",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_CustomReturnType.Name)

	case !isQuorumCallVariant && hasMethodOption(method, gorums.E_QfWithReq):
		// Only QC variants need to process replies.
		log.Fatalf(
			"%s.%s: cannot combine non-quorum call method with the '%s' option",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_QfWithReq.Name)

	case !hasMethodOption(method, gorums.E_Multicast) && method.Desc.IsStreamingClient():
		log.Fatalf(
			"%s.%s: client-server streams is only valid with the '%s' option",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_Multicast.Name)

	case hasMethodOption(method, gorums.E_Multicast) && !method.Desc.IsStreamingClient():
		log.Fatalf(
			"%s.%s: '%s' option is only valid for client-server streams methods",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_Multicast.Name)

	case !hasMethodOption(method, gorums.E_CorrectableStream) && method.Desc.IsStreamingServer():
		log.Fatalf(
			"%s.%s: server-client streams is only valid with the '%s' option",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_CorrectableStream.Name)

	case hasMethodOption(method, gorums.E_CorrectableStream) && !method.Desc.IsStreamingServer():
		log.Fatalf(
			"%s.%s: '%s' option is only valid for server-client streams",
			method.Parent.Desc.Name(), method.Desc.Name(), gorums.E_CorrectableStream.Name)
	}

	return firstOption
}
