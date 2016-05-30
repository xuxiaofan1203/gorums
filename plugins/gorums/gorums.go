// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2015 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package gorums outputs a gorums client API in Go code.
// It runs as a plugin for the Go protocol buffer compiler plugin.
// It is linked in to protoc-gen-go.
package gorums

import (
	"log"
	"sort"
	"strings"

	pb "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

func init() {
	generator.RegisterPlugin(new(gorums))
}

// gorums is an implementation of the Go protocol buffer compiler's
// plugin architecture.  It generates bindings for gorums support.
type gorums struct {
	gen *generator.Generator
}

// Name returns the name of this plugin, "gorums".
func (g *gorums) Name() string {
	return "gorums"
}

// Init initializes the plugin.
func (g *gorums) Init(gen *generator.Generator) {
	g.gen = gen
}

// Given a type name defined in a .proto, return its object.
// Also record that we're using it, to guarantee the associated import.
func (g *gorums) objectNamed(name string) generator.Object {
	g.gen.RecordTypeUse(name)
	return g.gen.ObjectNamed(name)
}

// Given a type name defined in a .proto, return its name as we will print it.
func (g *gorums) typeName(str string) string {
	return g.gen.TypeName(g.objectNamed(str))
}

// P forwards to g.gen.P.
func (g *gorums) P(args ...interface{}) { g.gen.P(args...) }

// Generate generates code for the services in the given file.
func (g *gorums) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	smethods := g.generateServiceMethods(file.FileDescriptorProto.Service)

	g.generateMgrTypeRelated(smethods)
	g.generateGorumsWrapperForService(file, smethods)
	g.embedStaticResources()
}

// GenerateImports generates the import declaration for this file.
func (g *gorums) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	if len(staticImports) == 0 {
		return
	}

	sort.Strings(staticImports)
	g.P("import (")
	for _, simport := range staticImports {
		if ignore := ignoreImport[simport]; ignore {
			continue
		}
		g.P("\"", simport, "\"")
	}
	if !onlyClientStreamMethods(file.FileDescriptorProto.Service) {
		g.P()
		g.P("\"google.golang.org/grpc/codes\"")
	}
	g.P(")")
}

var ignoreImport = map[string]bool{
	"fmt":  true,
	"math": true,
	"golang.org/x/net/context": true,
	"google.golang.org/grpc":   true,
}

func onlyClientStreamMethods(services []*pb.ServiceDescriptorProto) bool {
	for _, service := range services {
		for _, method := range service.Method {
			if !method.GetClientStreaming() {
				return false
			}
		}

	}
	return true
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func (g *gorums) embedStaticResources() {
	g.P("/* Static resources */")
	g.P(staticResources)
}

type serviceMethod struct {
	origName string
	name     string
	unexName string
	rpcName  string

	respName string
	requName string

	typeName     string
	unexTypeName string

	streaming bool

	servName string // Redundant, but keeps it simple.
}

type smSlice []serviceMethod

func (p smSlice) Len() int      { return len(p) }
func (p smSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p smSlice) Less(i, j int) bool {
	if p[i].servName < p[j].servName {
		return true
	} else if p[i].servName > p[j].servName {
		return false
	} else {
		if p[i].origName < p[j].origName {
			return true
		} else if p[i].origName > p[j].origName {
			return false
		} else {
			return false
		}
	}
}

func (g *gorums) generateServiceMethods(services []*pb.ServiceDescriptorProto) []serviceMethod {
	smethods := make(map[string][]*serviceMethod)
	for _, service := range services {
		for _, method := range service.Method {
			if method.GetServerStreaming() {
				log.Fatalf(
					"%s - %s: server streaming not supported by gorums",
					service.GetName(), method.GetName())
			}

			sm := serviceMethod{}
			sm.origName = method.GetName()
			sm.name = generator.CamelCase(sm.origName)
			sm.rpcName = sm.name // sm.Name may be overwritten if method name conflict
			sm.unexName = unexport(sm.name)
			sm.respName = g.typeName(method.GetOutputType())
			sm.requName = g.typeName(method.GetInputType())
			sm.typeName = sm.name + "Reply"
			sm.unexTypeName = unexport(sm.typeName)
			sm.servName = service.GetName()

			if sm.typeName == sm.respName {
				sm.typeName += "_"
			}

			if method.GetClientStreaming() {
				sm.streaming = true
			}

			methodsForName, _ := smethods[sm.name]
			methodsForName = append(methodsForName, &sm)
			smethods[sm.name] = methodsForName
		}
	}

	var allRewrittenFlat []serviceMethod

	for _, methodsForName := range smethods {
		switch len(methodsForName) {
		case 0:
			panic("generateServiceMethods: found method name with no data")
		case 1:
			allRewrittenFlat = append(allRewrittenFlat, *methodsForName[0])
			continue
		default:
			for _, sm := range methodsForName {
				sm.origName = sm.servName + sm.origName
				sm.name = sm.servName + sm.name
				sm.unexName = unexport(sm.name)
				sm.typeName = sm.name + "Reply"
				sm.unexTypeName = unexport(sm.typeName)
				if sm.typeName == sm.respName {
					sm.typeName += "_"
				}
				allRewrittenFlat = append(allRewrittenFlat, *sm)
			}
		}
	}

	sort.Sort(smSlice(allRewrittenFlat))

	return allRewrittenFlat
}

// generateGorumsWrapperForService generates all the code for the named service.
func (g *gorums) generateGorumsWrapperForService(file *generator.FileDescriptor, smethods []serviceMethod) {
	g.P()
	g.P("/* Gorums Client API */")
	g.P()
	g.generateRPCConfig(smethods)
	g.generateRPCMgr(file.GetPackage(), smethods)
}

func (g *gorums) generateRPCConfig(smethods []serviceMethod) {
	g.P()
	g.P("/* Configuration RPC specific */")
	for _, method := range smethods {
		g.P()
		if method.streaming {
			g.generateConfigStreamMethod(method)
		} else {
			g.generateConfigReplyTypeAndMethod(method)
		}
	}
}

func (g *gorums) generateRPCMgr(pkgName string, smethods []serviceMethod) {
	g.P()
	g.P("/* Manager RPC specific */")
	for _, method := range smethods {
		if method.streaming {
			g.generateMgrStreamMethod(method)
		} else {
			g.generateMgrReplyTypeAndMethod(pkgName, method)
		}
	}
}

func (g *gorums) generateConfigReplyTypeAndMethod(sm serviceMethod) {
	g.P("// ", sm.typeName, " encapsulates the reply from a ", sm.name, " RPC invocation.")
	g.P("// It contains the id of each node in the quorum that replied and a single")
	g.P("// reply.")
	g.P("type ", sm.typeName, " struct {")
	g.P("NodeIDs []int")
	g.P("Reply *", sm.respName)
	g.P("}")
	g.P()

	g.P("func (r ", sm.typeName, ") String() string {")
	g.P("return fmt.Sprintf(\"node ids: %v | answer: %v\", r.NodeIDs, r.Reply)")
	g.P("}")
	g.P()

	g.P("// ", sm.typeName, " invokes a ", sm.name, " RPC on configuration c")
	g.P("// and returns the result as a ", sm.typeName, ".")
	g.P("func (c *Configuration) ", sm.name, "(args *", sm.requName, ") (*", sm.typeName, ", error) {")
	g.P("return c.mgr.", sm.unexName, "(c.id, args)")
	g.P("}")

	g.P("// ", sm.name, "Future is a reference to an asynchronous ", sm.name, " RPC invocation.")
	g.P("type ", sm.name, "Future struct {")
	g.P("	reply *", sm.typeName)
	g.P("	err   error")
	g.P("	c     chan struct{}")
	g.P("}")

	g.P("// ", sm.name, "Future asynchronously invokes a ", sm.name, " RPC on configuration c and")
	g.P("// returns a ", sm.name, "Future which can be used to inspect the RPC reply and error")
	g.P("// when available.")
	g.P("func (c *Configuration) ", sm.name, "Future(args *", sm.requName, ") *", sm.name, "Future {")
	g.P("	f := new(", sm.name, "Future)")
	g.P("	f.c = make(chan struct{}, 1)")
	g.P("	go func() {")
	g.P("		defer close(f.c)")
	g.P("		f.reply, f.err = c.mgr.", sm.unexName, "(c.id, args)")
	g.P("	}()")
	g.P("	return f")
	g.P("}")

	g.P("// Get returns the reply and any error associated with the ", sm.name, "Future.")
	g.P("// The method blocks until a reply or error is available.")
	g.P("func (f *", sm.name, "Future) Get() (*", sm.typeName, ", error) {")
	g.P("	<-f.c")
	g.P("	return f.reply, f.err")
	g.P("}")

	g.P("// Done reports if a reply or error is available for the ", sm.name, "Future.")
	g.P("func (f *", sm.name, "Future) Done() bool {")
	g.P("select {")
	g.P("case <-f.c:")
	g.P("return true")
	g.P("default:")
	g.P("return false")
	g.P("}")
	g.P("}")
}

func (g *gorums) generateConfigStreamMethod(sm serviceMethod) {
	g.P()
	g.P("// ", sm.name, "invokes an asynchronous ", sm.name, " RPC on configuration c.")
	g.P("// The call has no return value and is invoked on every node in the")
	g.P("// configuration.")
	g.P("func (c *Configuration) ", sm.name, "(args *", sm.requName, ") error {")
	g.P("return c.mgr.", sm.unexName, "(c.id, args)")
	g.P("}")
}

func (g *gorums) generateMgrReplyTypeAndMethod(pkgName string, sm serviceMethod) {
	g.P()
	g.P("type ", sm.unexTypeName, " struct {")
	g.P("nid int")
	g.P("reply *", sm.respName)
	g.P("err error")
	g.P("}")

	g.P()
	g.P("func (m *Manager) ", sm.unexName, "(cid int, args *", sm.requName, ") (*", sm.typeName, ", error) {")
	g.P("c, found := m.Configuration(cid)")
	g.P("if !found {")
	g.P("panic(\"execptional: config not found\")")
	g.P("}")
	g.P()

	g.P("var (")
	g.P("replyChan = make(chan ", sm.unexTypeName, ", c.Size())")
	g.P("stopSignal  = make(chan struct{})")
	g.P("replyValues = make([]*", sm.respName, ", 0, c.quorum)")
	g.P("errCount int")
	g.P("quorum bool")
	g.P("reply = &", sm.typeName, "{NodeIDs: make([]int, 0, c.quorum)}")
	g.P("ctx, cancel = context.WithCancel(context.Background())")
	g.P(")")
	g.P()

	g.P("for _, nid := range c.nodes {")
	g.P("node, found := m.Node(nid)")
	g.P("if !found {")
	g.P("panic(\"exceptional: node not found\")")
	g.P("}")
	g.P("go func() {")
	g.P("reply := new(", sm.respName, ")")
	g.P("ce := make(chan error, 1)")
	g.P("start := time.Now()")
	g.P("go func() {")
	g.P("select {")
	g.P("case ce <- grpc.Invoke(")
	g.P("ctx,")
	g.P("\"/", pkgName, ".", sm.servName, "/", sm.rpcName, "\",")
	g.P("args,")
	g.P("reply,")
	g.P("node.conn,")
	g.P("):")
	g.P("case <-stopSignal:")
	g.P("return")
	g.P("}")
	g.P("}()")
	g.P("select {")
	g.P("case err := <-ce:")
	g.P("switch grpc.Code(err) {")
	g.P("case codes.OK, codes.Aborted, codes.Canceled:")
	g.P("node.setLatency(time.Since(start))")
	g.P("default:")
	g.P("node.setLastErr(err)")
	g.P("}")
	g.P("replyChan <- ", sm.unexTypeName, "{node.id, reply, err}")
	g.P("case <-stopSignal:")
	g.P("return")
	g.P("}")
	g.P("}()")
	g.P("}")
	g.P()

	g.P("defer close(stopSignal)")
	g.P("defer cancel()")
	g.P()

	g.P("for {")
	g.P()
	g.P("select {")
	g.P("case r := <-replyChan:")
	g.P("if r.err != nil {")
	g.P("errCount++")
	g.P("goto terminationCheck")
	g.P("}")
	g.P()
	g.P("replyValues = append(replyValues, r.reply)")
	g.P("reply.NodeIDs = append(reply.NodeIDs, r.nid)")
	g.P("if reply.Reply, quorum = m.", sm.unexName, "qf(c, replyValues); quorum {")
	g.P("return reply, nil")
	g.P("}")
	g.P("case <-time.After(c.timeout):")
	g.P("return reply, TimeoutRPCError{c.timeout, errCount, len(replyValues)}")
	g.P("}")
	g.P()

	g.P("terminationCheck:")
	g.P("if errCount+len(replyValues) == c.Size() {")
	g.P("return reply, IncompleteRPCError{errCount, len(replyValues)}")
	g.P("}")
	g.P("}")
	g.P("}")
}

func (g *gorums) generateMgrStreamMethod(sm serviceMethod) {
	g.P("func (m *Manager) ", sm.unexName, "(cid int, args *", sm.requName, ") error {")
	g.P("c, found := m.Configuration(cid)")
	g.P("if !found {")
	g.P("panic(\"execeptional: config not found\")")
	g.P("}")
	g.P()
	g.P("for _, nid := range c.nodes {")
	g.P("go func(nodeID int) {")
	g.P("stream := m.", sm.unexName, "Clients[nodeID]")
	g.P("if stream == nil {")
	g.P("panic(\"execeptional: node client stream not found\")")
	g.P("}")
	g.P("err := stream.Send(args)")
	g.P("if err == nil {")
	g.P("return")
	g.P("}")
	g.P("if m.logger != nil {")
	g.P("m.logger.Printf(\"node %d: ", sm.unexName, " stream send error: %v\", nodeID, err)")
	g.P("}")
	g.P("}(nid)")
	g.P("}")
	g.P()
	g.P("return nil")
	g.P("}")
}

func (g *gorums) generateMgrTypeRelated(smethods []serviceMethod) {
	g.P("/* Manager type struct */")
	g.generateMgrStruct(smethods)

	g.P("/* Manager quorum functions */")
	g.generateMgrSetDefaultQFuncs(smethods)

	g.P("/* Manager create/close streams */")
	g.generateMgrCreateStreams(smethods)
	g.generateMgrCloseStreams(smethods)

	g.P("/* Manager options */")
	g.generateMgrOptions(smethods)

	for _, method := range smethods {
		if method.streaming {
			continue
		}
		g.P()
		g.generateQuorumFunctionType(method)
	}
}

func (g *gorums) generateMgrStruct(sms []serviceMethod) {
	g.P()
	g.P("// Manager manages a pool of node configurations on which quorum remote")
	g.P("// procedure calls can be made.")
	g.P("type Manager struct {")
	g.P("sync.RWMutex")
	g.P("nodes []*Node")
	g.P("configs []*Configuration")
	g.P("nodeGidToID map[uint32]int")
	g.P("configGidToID map[uint32]int")
	g.P()
	g.P("closeOnce sync.Once")
	g.P("logger *log.Logger")
	g.P("opts managerOptions")
	g.P()
	for _, sm := range sms {
		if sm.streaming {
			continue
		}
		g.P(sm.unexName, "qf ", sm.name, "QuorumFn")
	}
	g.P()
	for _, sm := range sms {
		if !sm.streaming {
			continue
		}
		g.P(sm.unexName, "Clients []", sm.servName, "_", sm.name, "Client")
	}
	g.P("}")
}

func (g *gorums) generateMgrSetDefaultQFuncs(sms []serviceMethod) {
	g.P()
	g.P("func (m *Manager) setDefaultQuorumFuncs() {")
	for _, sm := range sms {
		if sm.streaming {
			continue
		}
		g.P("if m.opts.", sm.unexName, "qf != nil {")
		g.P("m.", sm.unexName, "qf = m.opts.", sm.unexName, "qf")
		g.P("} else {")
		g.P("m.", sm.unexName, "qf = func(c *Configuration, replies []*", sm.respName, ") (*", sm.respName, ", bool) {")
		g.P("if len(replies) < c.Quorum() {")
		g.P("return nil, false")
		g.P("}")
		g.P("return replies[0], true")
		g.P("}")
		g.P("}")
	}
	g.P("}")
}

func (g *gorums) generateMgrCreateStreams(sms []serviceMethod) {
	g.P()
	g.P("func (m *Manager) createStreamClients() error {")
	g.P("if m.opts.noConnect {")
	g.P("return nil")
	g.P("}")
	g.P()
	for _, sm := range sms {
		if !sm.streaming {
			continue
		}
		g.P("for _, node := range m.nodes {")
		g.P("if node.self {")
		g.P("m.", sm.unexName, "Clients = append(m.", sm.unexName, "Clients, nil)")
		g.P("continue")
		g.P("}")
		g.P("client := New", sm.servName, "Client(node.conn)")
		g.P(sm.unexName, "Client, err := client.", sm.name, "(context.Background())")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("m.", sm.unexName, "Clients = append(m.", sm.unexName, "Clients, ", sm.unexName, "Client)")
		g.P("}")
	}
	g.P()
	g.P("return nil")
	g.P("}")
}

func (g *gorums) generateMgrCloseStreams(sms []serviceMethod) {
	g.P()
	g.P("func (m *Manager) closeStreamClients() {")
	g.P("if m.opts.noConnect {")
	g.P("return")
	g.P("}")
	g.P()
	for _, sm := range sms {
		if !sm.streaming {
			continue
		}
		g.P("for i, client := range m.", sm.unexName, "Clients {")
		g.P("_, err := client.CloseAndRecv()")
		g.P("if err == nil {")
		g.P("continue")
		g.P("}")
		g.P("if m.logger != nil {")
		g.P("m.logger.Printf(\"node %d: error closing ", sm.unexName, " client: %v\", i, err)")
		g.P("}")
		g.P("}")
	}
	g.P("}")
}

func (g *gorums) generateQuorumFunctionType(sm serviceMethod) {
	g.P()
	g.P("// ", sm.name, "QuorumFn is used to pick a reply from the replies if there is a quorum.")
	g.P("// If there was not enough replies to satisfy the quorum requirement,")
	g.P("// then the function returns (nil, false). Otherwise, the function picks a")
	g.P("// reply among the replies and returns (reply, true).")
	g.P("type ", sm.name, "QuorumFn func(c *Configuration, replies []*", sm.respName, ") (*", sm.respName, ", bool)")
}

func (g *gorums) generateMgrOptions(sms []serviceMethod) {
	g.P()
	g.P("type managerOptions struct {")
	g.P("grpcDialOpts []grpc.DialOption")
	g.P("logger *log.Logger")
	g.P("noConnect bool")
	g.P("selfAddr string")
	g.P("selfGid uint32")
	g.P()
	for _, sm := range sms {
		if sm.streaming {
			continue
		}
		g.P(sm.unexName, "qf ", sm.name, "QuorumFn")
	}
	g.P("}")
	g.P()
	for _, sm := range sms {
		if sm.streaming {
			continue
		}
		g.P("// With", sm.name, "QuorumFunc returns a ManagerOption that sets a cumstom")
		g.P("// ", sm.name, "QuorumFunc.")
		g.P("func With", sm.name, "QuorumFunc(f ", sm.name, "QuorumFn) ManagerOption {")
		g.P("return func(o *managerOptions) {")
		g.P("o.", sm.unexName, "qf = f")
		g.P("}")
		g.P("}")
	}
}
