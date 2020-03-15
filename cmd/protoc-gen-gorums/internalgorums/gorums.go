// Package internalgorums is internal to the gorums protobuf module.
package internalgorums

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/relab/gorums"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// TODO(meling) replace github.com/relab/gorums with gorums.io as import package

var importMap = map[string]protogen.GoImportPath{
	"io":      protogen.GoImportPath("io"),
	"time":    protogen.GoImportPath("time"),
	"fmt":     protogen.GoImportPath("fmt"),
	"log":     protogen.GoImportPath("log"),
	"sync":    protogen.GoImportPath("sync"),
	"context": protogen.GoImportPath("context"),
	"trace":   protogen.GoImportPath("golang.org/x/net/trace"),
	"grpc":    protogen.GoImportPath("google.golang.org/grpc"),
	"codes":   protogen.GoImportPath("google.golang.org/grpc/codes"),
	"status":  protogen.GoImportPath("google.golang.org/grpc/status"),
	"gorums":  protogen.GoImportPath("github.com/relab/gorums"),
}

func addImport(path, ident string, g *protogen.GeneratedFile) string {
	pkg := path[strings.LastIndex(path, "/")+1:]
	impPath, ok := importMap[pkg]
	if !ok {
		impPath = protogen.GoImportPath(path)
		importMap[pkg] = impPath
	}
	return g.QualifiedGoIdent(impPath.Ident(ident))
}

type servicesData struct {
	GenFile  *protogen.GeneratedFile
	Services []*protogen.Service
}

type methodData struct {
	GenFile *protogen.GeneratedFile
	Method  *protogen.Method
}

// GenerateFile generates a _gorums.pb.go file containing Gorums service definitions.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 || !checkMethodOptions(file.Services, gorumsCallTypes...) {
		// there is nothing for this plugin to do
		return nil
	}
	if len(file.Services) > 1 {
		// To build multiple services, make separate proto files and
		// run the plugin separately for each proto file.
		// These cannot share the same Go package.
		log.Fatalln("Gorums does not support multiple services in the same proto file.")
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
	g.P(mustExecute(parseTemplate("Node", node), data))
	g.P()
	g.P(mustExecute(parseTemplate("QuorumSpec", qspecInterface), data))
	g.P()
	g.P(mustExecute(parseTemplate("DataTypes", datatypes), data))
	g.P()
	genGorumsMethods(servicesData{g, file.Services}, gorumsCallTypes...)
	g.P()
}

func genGorumsMethods(data servicesData, methodOptions ...*protoimpl.ExtensionInfo) {
	g := data.GenFile
	for _, service := range data.Services {
		for _, method := range service.Methods {
			if hasMethodOption(method, methodOptions...) {
				fmt.Fprintf(os.Stderr, "processing %s\n", method.GoName)
				g.P(genGorumsMethod(g, method))
			}
		}
	}
}

func genGorumsMethod(g *protogen.GeneratedFile, method *protogen.Method) string {
	methodOption := validateMethodExtensions(method)
	if template, ok := gorumsCallTypeTemplates[methodOption]; ok {
		return mustExecute(parseTemplate(methodOption.Name, template), methodData{g, method})
	}
	panic(fmt.Sprintf("unknown method type %s\n", method.GoName))
}

// hasGorumsType returns true if one of the service methods specify
// the given gorums type.
func hasGorumsType(services []*protogen.Service, gorumsType string) bool {
	// TODO(meling) try to avoid this loop slice; reuse devTypes??
	for _, gType := range []string{"node", "qspec", "types"} {
		if gorumsType == gType {
			return true
		}
	}
	if methodOption, ok := gorumsTypes[gorumsType]; ok {
		return checkMethodOptions(services, methodOption)
	}
	return false
}

// compute index to start of option name
const index = len("gorums.")

// name to method option mapping
var gorumsTypes = map[string]*protoimpl.ExtensionInfo{
	gorums.E_Qc.Name[index:]:                gorums.E_Qc,
	gorums.E_QcFuture.Name[index:]:          gorums.E_QcFuture,
	gorums.E_Correctable.Name[index:]:       gorums.E_Correctable,
	gorums.E_CorrectableStream.Name[index:]: gorums.E_CorrectableStream,
	gorums.E_Multicast.Name[index:]:         gorums.E_Multicast,
}

var gorumsCallTypeTemplates = map[*protoimpl.ExtensionInfo]string{
	gorums.E_Qc:                quorumCall,
	gorums.E_QcFuture:          futureCall,
	gorums.E_Correctable:       correctableCall,
	gorums.E_CorrectableStream: correctableStreamCall,
	gorums.E_Multicast:         multicastCall,
}

// gorumsCallTypes should list all available call types supported by Gorums.
// These are considered mutually incompatible.
var gorumsCallTypes = []*protoimpl.ExtensionInfo{
	gorums.E_Qc,
	gorums.E_QcFuture,
	gorums.E_Correctable,
	gorums.E_CorrectableStream,
	gorums.E_Multicast,
}

// callTypesWithInternal should list all available call types that
// has a quorum function and hence need an internal type that wraps
// the return type with additional information.
var callTypesWithInternal = []*protoimpl.ExtensionInfo{
	gorums.E_Qc,
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

// hasGorumsCallType returns true if the given method has specified
// one of the call types supported by Gorums.
func hasGorumsCallType(method *protogen.Method) bool {
	return hasMethodOption(method, gorumsCallTypes...)
}

// checkMethodOptions returns true if one of the methods provided by
// the given services has one of the given options.
func checkMethodOptions(services []*protogen.Service, methodOptions ...*protoimpl.ExtensionInfo) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			if hasMethodOption(method, methodOptions...) {
				return true
			}
		}
	}
	return false
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
