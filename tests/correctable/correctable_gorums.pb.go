// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.6.0-devel
// 	protoc            v3.17.3
// source: correctable/correctable.proto

package correctable

import (
	context "context"
	fmt "fmt"
	gorums "github.com/relab/gorums"
	ordering "github.com/relab/gorums/ordering"
	encoding "google.golang.org/grpc/encoding"
	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = gorums.EnforceVersion(6 - gorums.MinVersion)
	// Verify that the gorums runtime is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.MaxVersion - 6)
)

// A Configuration represents a static set of nodes on which quorum remote
// procedure calls may be invoked.
type Configuration struct {
	gorums.Configuration
	qspec QuorumSpec
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Manager.
func (c *Configuration) Nodes() []*Node {
	nodes := make([]*Node, 0, c.Size())
	for _, n := range c.Configuration {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

// And returns a NodeListOption that can be used to create a new configuration combining c and d.
func (c Configuration) And(d *Configuration) gorums.NodeListOption {
	return c.Configuration.And(d.Configuration)
}

// Except returns a NodeListOption that can be used to create a new configuration
// from c without the nodes in rm.
func (c Configuration) Except(rm *Configuration) gorums.NodeListOption {
	return c.Configuration.Except(rm.Configuration)
}

func init() {
	if encoding.GetCodec(gorums.ContentSubtype) == nil {
		encoding.RegisterCodec(gorums.NewCodec())
	}
}

// Manager maintains a connection pool of nodes on
// which quorum calls can be performed.
type Manager struct {
	*gorums.Manager
}

// NewManager returns a new Manager for managing connection to nodes added
// to the manager. This function accepts manager options used to configure
// various aspects of the manager.
func NewManager(opts ...gorums.ManagerOption) (mgr *Manager) {
	mgr = &Manager{}
	mgr.Manager = gorums.NewManager(opts...)
	return mgr
}

// NewConfiguration returns a configuration based on the provided list of nodes (required)
// and an optional quorum specification. The QuorumSpec is necessary for call types that
// must process replies. For configurations only used for unicast or multicast call types,
// a QuorumSpec is not needed. The QuorumSpec interface is also a ConfigOption.
// Nodes can be supplied using WithNodeMap or WithNodeList, or WithNodeIDs.
// A new configuration can also be created from an existing configuration,
// using the And, WithNewNodes, Except, and WithoutNodes methods.
func (m *Manager) NewConfiguration(opts ...gorums.ConfigOption) (c *Configuration, err error) {
	if len(opts) < 1 || len(opts) > 2 {
		return nil, fmt.Errorf("wrong number of options: %d", len(opts))
	}
	c = &Configuration{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case gorums.NodeListOption:
			c.Configuration, err = gorums.NewConfiguration(m.Manager, v)
			if err != nil {
				return nil, err
			}
		case QuorumSpec:
			// Must be last since v may match QuorumSpec if it is interface{}
			c.qspec = v
		default:
			return nil, fmt.Errorf("unknown option type: %v", v)
		}
	}
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && c.qspec == nil {
		return nil, fmt.Errorf("missing required QuorumSpec")
	}
	return c, nil
}

// Nodes returns a slice of available nodes on this manager.
// IDs are returned in the order they were added at creation of the manager.
func (m *Manager) Nodes() []*Node {
	gorumsNodes := m.Manager.Nodes()
	nodes := make([]*Node, 0, len(gorumsNodes))
	for _, n := range gorumsNodes {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

// Node encapsulates the state of a node on which a remote procedure call
// can be performed.
type Node struct {
	*gorums.Node
}

// Correctable asynchronously invokes a correctable quorum call on each node
// in configuration c and returns a CorrectableCorrectableResponse, which can be used
// to inspect any replies or errors when available.
func (c *Configuration) Correctable(ctx context.Context, in *CorrectableRequest) *CorrectableCorrectableResponse {
	cd := gorums.CorrectableCallData{
		Message:      in,
		Method:       "correctable.CorrectableTest.Correctable",
		ServerStream: false,
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, int, bool) {
		r := make(map[uint32]*CorrectableResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*CorrectableResponse)
		}
		return c.qspec.CorrectableQF(req.(*CorrectableRequest), r)
	}

	corr := c.Configuration.CorrectableCall(ctx, cd)
	return &CorrectableCorrectableResponse{corr}
}

// CorrectableStream asynchronously invokes a correctable quorum call on each node
// in configuration c and returns a CorrectableStreamCorrectableResponse, which can be used
// to inspect any replies or errors when available.
// This method supports server-side preliminary replies (correctable stream).
func (c *Configuration) CorrectableStream(ctx context.Context, in *CorrectableRequest) *CorrectableStreamCorrectableResponse {
	cd := gorums.CorrectableCallData{
		Message:      in,
		Method:       "correctable.CorrectableTest.CorrectableStream",
		ServerStream: true,
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, int, bool) {
		r := make(map[uint32]*CorrectableResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*CorrectableResponse)
		}
		return c.qspec.CorrectableStreamQF(req.(*CorrectableRequest), r)
	}

	corr := c.Configuration.CorrectableCall(ctx, cd)
	return &CorrectableStreamCorrectableResponse{corr}
}

// QuorumSpec is the interface of quorum functions for CorrectableTest.
type QuorumSpec interface {
	gorums.ConfigOption

	// CorrectableQF is the quorum function for the Correctable
	// correctable quorum call method. The in parameter is the request object
	// supplied to the Correctable method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *CorrectableRequest'.
	CorrectableQF(in *CorrectableRequest, replies map[uint32]*CorrectableResponse) (*CorrectableResponse, int, bool)

	// CorrectableStreamQF is the quorum function for the CorrectableStream
	// correctable stream quorum call method. The in parameter is the request object
	// supplied to the CorrectableStream method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *CorrectableRequest'.
	CorrectableStreamQF(in *CorrectableRequest, replies map[uint32]*CorrectableResponse) (*CorrectableResponse, int, bool)
}

// CorrectableTest is the server-side API for the CorrectableTest Service
type CorrectableTest interface {
	Correctable(ctx gorums.ServerCtx, request *CorrectableRequest) (response *CorrectableResponse, err error)
	CorrectableStream(ctx gorums.ServerCtx, request *CorrectableRequest, send func(response *CorrectableResponse) error) error
}

func RegisterCorrectableTestServer(srv *gorums.Server, impl CorrectableTest) {
	srv.RegisterHandler("correctable.CorrectableTest.Correctable", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*CorrectableRequest)
		defer ctx.Release()
		resp, err := impl.Correctable(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("correctable.CorrectableTest.CorrectableStream", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*CorrectableRequest)
		defer ctx.Release()
		err := impl.CorrectableStream(ctx, req, func(resp *CorrectableResponse) error {
			// create a copy of the metadata, to avoid a data race between WrapMessage and SendMsg
			md := proto.Clone(in.Metadata)
			return gorums.SendMessage(ctx, finished, gorums.WrapMessage(md.(*ordering.Metadata), resp, nil))
		})
		if err != nil {
			gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, nil, err))
		}
	})
}

type internalCorrectableResponse struct {
	nid   uint32
	reply *CorrectableResponse
	err   error
}

// CorrectableCorrectableResponse is a correctable object for processing replies.
type CorrectableCorrectableResponse struct {
	*gorums.Correctable
}

// Get returns the reply, level and any error associated with the
// called method. The method does not block until a (possibly
// intermediate) reply or error is available. Level is set to LevelNotSet if no
// reply has yet been received. The Done or Watch methods should be used to
// ensure that a reply is available.
func (c *CorrectableCorrectableResponse) Get() (*CorrectableResponse, int, error) {
	resp, level, err := c.Correctable.Get()
	if err != nil {
		return nil, level, err
	}
	return resp.(*CorrectableResponse), level, err
}

// CorrectableStreamCorrectableResponse is a correctable object for processing replies.
type CorrectableStreamCorrectableResponse struct {
	*gorums.Correctable
}

// Get returns the reply, level and any error associated with the
// called method. The method does not block until a (possibly
// intermediate) reply or error is available. Level is set to LevelNotSet if no
// reply has yet been received. The Done or Watch methods should be used to
// ensure that a reply is available.
func (c *CorrectableStreamCorrectableResponse) Get() (*CorrectableResponse, int, error) {
	resp, level, err := c.Correctable.Get()
	if err != nil {
		return nil, level, err
	}
	return resp.(*CorrectableResponse), level, err
}
