// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	fmt "fmt"
	ptypes "github.com/golang/protobuf/ptypes"
	gorums "github.com/relab/gorums"
	trace "golang.org/x/net/trace"
	grpc "google.golang.org/grpc"
	time "time"
)

// ---------------------------------------------------------------
// Strict Ordering RPCs
// ---------------------------------------------------------------
func (c *Configuration) StrictOrderingQC(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "StrictOrderingQC")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *strictOrderingResult, c.n)
	c.mgr.recvQMut.Lock()
	c.mgr.recvQ[msgID] = replies
	c.mgr.recvQMut.Unlock()

	defer func() {
		// remove the replies channel when we are done
		c.mgr.recvQMut.Lock()
		delete(c.mgr.recvQ, msgID)
		c.mgr.recvQMut.Unlock()
	}()

	data, err := ptypes.MarshalAny(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &gorums.Message{
		ID:   msgID,
		URL:  "/dev.ZorumsService/StrictOrderingQC",
		Data: data,
	}
	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		n.strictOrdering.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := ptypes.UnmarshalAny(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.StrictOrderingQCQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}

		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

// StrictOrderingQCHandler is the server API for the StrictOrderingQC rpc.
type StrictOrderingQCHandler interface {
	StrictOrderingQC(*Request) *Response
}

// RegisterStrictOrderingQCHandler sets the handler for StrictOrderingQC.
func (s *GorumsServer) RegisterStrictOrderingQCHandler(handler StrictOrderingQCHandler) {
	s.srv.registerHandler("/dev.ZorumsService/StrictOrderingQC", func(in *gorums.Message) *gorums.Message {
		req := new(Request)
		err := ptypes.UnmarshalAny(in.GetData(), req)
		// TODO: how to handle marshaling errors here
		if err != nil {
			return new(gorums.Message)
		}
		resp := handler.StrictOrderingQC(req)
		data, err := ptypes.MarshalAny(resp)
		if err != nil {
			return new(gorums.Message)
		}
		return &gorums.Message{Data: data, URL: in.GetURL()}
	})
}

// StrictOrderingPerNodeArg is a quorum call invoked on each node in configuration c,
// with the argument returned by the provided function f, and returns the combined result.
// The per node function f receives a copy of the Request request argument and
// returns a Request manipulated to be passed to the given nodeID.
// The function f must be thread-safe.
func (c *Configuration) StrictOrderingPerNodeArg(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "StrictOrderingPerNodeArg")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *strictOrderingResult, c.n)
	c.mgr.recvQMut.Lock()
	c.mgr.recvQ[msgID] = replies
	c.mgr.recvQMut.Unlock()

	defer func() {
		// remove the replies channel when we are done
		c.mgr.recvQMut.Lock()
		delete(c.mgr.recvQ, msgID)
		c.mgr.recvQMut.Unlock()
	}()

	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		nodeArg := f(in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		data, err := ptypes.MarshalAny(nodeArg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}
		msg := &gorums.Message{
			ID:   msgID,
			URL:  "/dev.ZorumsService/StrictOrderingPerNodeArg",
			Data: data,
		}
		n.strictOrdering.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := ptypes.UnmarshalAny(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.StrictOrderingPerNodeArgQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}

		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

// StrictOrderingPerNodeArgHandler is the server API for the StrictOrderingPerNodeArg rpc.
type StrictOrderingPerNodeArgHandler interface {
	StrictOrderingPerNodeArg(*Request) *Response
}

// RegisterStrictOrderingPerNodeArgHandler sets the handler for StrictOrderingPerNodeArg.
func (s *GorumsServer) RegisterStrictOrderingPerNodeArgHandler(handler StrictOrderingPerNodeArgHandler) {
	s.srv.registerHandler("/dev.ZorumsService/StrictOrderingPerNodeArg", func(in *gorums.Message) *gorums.Message {
		req := new(Request)
		err := ptypes.UnmarshalAny(in.GetData(), req)
		// TODO: how to handle marshaling errors here
		if err != nil {
			return new(gorums.Message)
		}
		resp := handler.StrictOrderingPerNodeArg(req)
		data, err := ptypes.MarshalAny(resp)
		if err != nil {
			return new(gorums.Message)
		}
		return &gorums.Message{Data: data, URL: in.GetURL()}
	})
}

// StrictOrderingQFWithReq is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StrictOrderingQFWithReq(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *Response, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "StrictOrderingQFWithReq")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *strictOrderingResult, c.n)
	c.mgr.recvQMut.Lock()
	c.mgr.recvQ[msgID] = replies
	c.mgr.recvQMut.Unlock()

	defer func() {
		// remove the replies channel when we are done
		c.mgr.recvQMut.Lock()
		delete(c.mgr.recvQ, msgID)
		c.mgr.recvQMut.Unlock()
	}()

	data, err := ptypes.MarshalAny(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &gorums.Message{
		ID:   msgID,
		URL:  "/dev.ZorumsService/StrictOrderingQFWithReq",
		Data: data,
	}
	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		n.strictOrdering.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := ptypes.UnmarshalAny(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.StrictOrderingQFWithReqQF(in, replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}

		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

// StrictOrderingQFWithReqHandler is the server API for the StrictOrderingQFWithReq rpc.
type StrictOrderingQFWithReqHandler interface {
	StrictOrderingQFWithReq(*Request) *Response
}

// RegisterStrictOrderingQFWithReqHandler sets the handler for StrictOrderingQFWithReq.
func (s *GorumsServer) RegisterStrictOrderingQFWithReqHandler(handler StrictOrderingQFWithReqHandler) {
	s.srv.registerHandler("/dev.ZorumsService/StrictOrderingQFWithReq", func(in *gorums.Message) *gorums.Message {
		req := new(Request)
		err := ptypes.UnmarshalAny(in.GetData(), req)
		// TODO: how to handle marshaling errors here
		if err != nil {
			return new(gorums.Message)
		}
		resp := handler.StrictOrderingQFWithReq(req)
		data, err := ptypes.MarshalAny(resp)
		if err != nil {
			return new(gorums.Message)
		}
		return &gorums.Message{Data: data, URL: in.GetURL()}
	})
}

// StrictOrderingCustomReturnType is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) StrictOrderingCustomReturnType(ctx context.Context, in *Request, opts ...grpc.CallOption) (resp *MyResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "StrictOrderingCustomReturnType")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *strictOrderingResult, c.n)
	c.mgr.recvQMut.Lock()
	c.mgr.recvQ[msgID] = replies
	c.mgr.recvQMut.Unlock()

	defer func() {
		// remove the replies channel when we are done
		c.mgr.recvQMut.Lock()
		delete(c.mgr.recvQ, msgID)
		c.mgr.recvQMut.Unlock()
	}()

	data, err := ptypes.MarshalAny(in)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	msg := &gorums.Message{
		ID:   msgID,
		URL:  "/dev.ZorumsService/StrictOrderingCustomReturnType",
		Data: data,
	}
	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		n.strictOrdering.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := ptypes.UnmarshalAny(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.StrictOrderingCustomReturnTypeQF(replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}

		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

// StrictOrderingCustomReturnTypeHandler is the server API for the StrictOrderingCustomReturnType rpc.
type StrictOrderingCustomReturnTypeHandler interface {
	StrictOrderingCustomReturnType(*Request) *Response
}

// RegisterStrictOrderingCustomReturnTypeHandler sets the handler for StrictOrderingCustomReturnType.
func (s *GorumsServer) RegisterStrictOrderingCustomReturnTypeHandler(handler StrictOrderingCustomReturnTypeHandler) {
	s.srv.registerHandler("/dev.ZorumsService/StrictOrderingCustomReturnType", func(in *gorums.Message) *gorums.Message {
		req := new(Request)
		err := ptypes.UnmarshalAny(in.GetData(), req)
		// TODO: how to handle marshaling errors here
		if err != nil {
			return new(gorums.Message)
		}
		resp := handler.StrictOrderingCustomReturnType(req)
		data, err := ptypes.MarshalAny(resp)
		if err != nil {
			return new(gorums.Message)
		}
		return &gorums.Message{Data: data, URL: in.GetURL()}
	})
}

// StrictOrderingCombi is a quorum call invoked on each node in configuration c,
// with the argument returned by the provided function f, and returns the combined result.
// The per node function f receives a copy of the Request request argument and
// returns a Request manipulated to be passed to the given nodeID.
// The function f must be thread-safe.
func (c *Configuration) StrictOrderingCombi(ctx context.Context, in *Request, f func(*Request, uint32) *Request, opts ...grpc.CallOption) (resp *MyResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "StrictOrderingCombi")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: in}, false)

		defer func() {
			ti.LazyLog(&qcresult{reply: resp, err: err}, false)
			if err != nil {
				ti.SetError()
			}
		}()
	}

	// get the ID which will be used to return the correct responses for a request
	msgID := c.mgr.nextMsgID()

	// set up a channel to collect replies
	replies := make(chan *strictOrderingResult, c.n)
	c.mgr.recvQMut.Lock()
	c.mgr.recvQ[msgID] = replies
	c.mgr.recvQMut.Unlock()

	defer func() {
		// remove the replies channel when we are done
		c.mgr.recvQMut.Lock()
		delete(c.mgr.recvQ, msgID)
		c.mgr.recvQMut.Unlock()
	}()

	// push the message to the nodes
	expected := c.n
	for _, n := range c.nodes {
		nodeArg := f(in, n.ID())
		if nodeArg == nil {
			expected--
			continue
		}
		data, err := ptypes.MarshalAny(nodeArg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal message: %w", err)
		}
		msg := &gorums.Message{
			ID:   msgID,
			URL:  "/dev.ZorumsService/StrictOrderingCombi",
			Data: data,
		}
		n.strictOrdering.sendQ <- msg
	}

	var (
		replyValues = make([]*Response, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replies:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			reply := new(Response)
			err := ptypes.UnmarshalAny(r.reply, reply)
			if err != nil {
				errs = append(errs, GRPCError{r.nid, fmt.Errorf("failed to unmarshal reply: %w", err)})
				break
			}
			replyValues = append(replyValues, reply)
			if resp, quorum = c.qspec.StrictOrderingCombiQF(in, replyValues); quorum {
				return resp, nil
			}
		case <-ctx.Done():
			return resp, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}
		}

		if len(errs)+len(replyValues) == expected {
			return resp, QuorumCallError{"incomplete call", len(replyValues), errs}
		}
	}
}

// StrictOrderingCombiHandler is the server API for the StrictOrderingCombi rpc.
type StrictOrderingCombiHandler interface {
	StrictOrderingCombi(*Request) *Response
}

// RegisterStrictOrderingCombiHandler sets the handler for StrictOrderingCombi.
func (s *GorumsServer) RegisterStrictOrderingCombiHandler(handler StrictOrderingCombiHandler) {
	s.srv.registerHandler("/dev.ZorumsService/StrictOrderingCombi", func(in *gorums.Message) *gorums.Message {
		req := new(Request)
		err := ptypes.UnmarshalAny(in.GetData(), req)
		// TODO: how to handle marshaling errors here
		if err != nil {
			return new(gorums.Message)
		}
		resp := handler.StrictOrderingCombi(req)
		data, err := ptypes.MarshalAny(resp)
		if err != nil {
			return new(gorums.Message)
		}
		return &gorums.Message{Data: data, URL: in.GetURL()}
	})
}
