// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	trace "golang.org/x/net/trace"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	time "time"
)

// ReadQuorumCall is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) ReadQuorumCall(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (resp *ReadResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadQuorumCall")
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

	expected := c.n
	replyChan := make(chan internalReadResponse, expected)
	for _, n := range c.nodes {
		go n.ReadQuorumCall(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*ReadResponse, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadQuorumCallQF(replyValues); quorum {
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

func (n *Node) ReadQuorumCall(ctx context.Context, in *ReadRequest, replyChan chan<- internalReadResponse) {
	reply := new(ReadResponse)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ReaderService/ReadQuorumCall", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalReadResponse{n.id, reply, err}
}

// ReadQuorumCallPerNodeArg is a quorum call invoked on each node in configuration c,
// with the argument returned by the provided function f, and returns the combined result.
// The per node function f receives a copy of the ReadRequest request argument and
// returns a ReadRequest manipulated to be passed to the given nodeID.
// The function f must be thread-safe.
func (c *Configuration) ReadQuorumCallPerNodeArg(ctx context.Context, in *ReadRequest, f func(*ReadRequest, uint32) *ReadRequest, opts ...grpc.CallOption) (resp *ReadResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadQuorumCallPerNodeArg")
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

	expected := c.n
	replyChan := make(chan internalReadResponse, expected)
	for _, n := range c.nodes {
		nodeArg := f(in, n.id)
		if nodeArg == nil {
			expected--
			continue
		}
		go n.ReadQuorumCallPerNodeArg(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*ReadResponse, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadQuorumCallPerNodeArgQF(replyValues); quorum {
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

func (n *Node) ReadQuorumCallPerNodeArg(ctx context.Context, in *ReadRequest, replyChan chan<- internalReadResponse) {
	reply := new(ReadResponse)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallPerNodeArg", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalReadResponse{n.id, reply, err}
}

// ReadQuorumCallQFWithRequestArg is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) ReadQuorumCallQFWithRequestArg(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (resp *ReadResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadQuorumCallQFWithRequestArg")
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

	expected := c.n
	replyChan := make(chan internalReadResponse, expected)
	for _, n := range c.nodes {
		go n.ReadQuorumCallQFWithRequestArg(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*ReadResponse, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadQuorumCallQFWithRequestArgQF(in, replyValues); quorum {
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

func (n *Node) ReadQuorumCallQFWithRequestArg(ctx context.Context, in *ReadRequest, replyChan chan<- internalReadResponse) {
	reply := new(ReadResponse)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallQFWithRequestArg", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalReadResponse{n.id, reply, err}
}

// ReadQuorumCallCustomReturnType is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) ReadQuorumCallCustomReturnType(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (resp *MyReadResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadQuorumCallCustomReturnType")
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

	expected := c.n
	replyChan := make(chan internalReadResponse, expected)
	for _, n := range c.nodes {
		go n.ReadQuorumCallCustomReturnType(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*ReadResponse, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadQuorumCallCustomReturnTypeQF(replyValues); quorum {
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

func (n *Node) ReadQuorumCallCustomReturnType(ctx context.Context, in *ReadRequest, replyChan chan<- internalReadResponse) {
	reply := new(ReadResponse)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallCustomReturnType", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalReadResponse{n.id, reply, err}
}

// ReadQuorumCallCombo does it all. Comment testing.
func (c *Configuration) ReadQuorumCallCombo(ctx context.Context, in *ReadRequest, f func(*ReadRequest, uint32) *ReadRequest, opts ...grpc.CallOption) (resp *MyReadResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadQuorumCallCombo")
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

	expected := c.n
	replyChan := make(chan internalReadResponse, expected)
	for _, n := range c.nodes {
		nodeArg := f(in, n.id)
		if nodeArg == nil {
			expected--
			continue
		}
		go n.ReadQuorumCallCombo(ctx, nodeArg, replyChan)
	}

	var (
		replyValues = make([]*ReadResponse, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadQuorumCallComboQF(in, replyValues); quorum {
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

func (n *Node) ReadQuorumCallCombo(ctx context.Context, in *ReadRequest, replyChan chan<- internalReadResponse) {
	reply := new(ReadResponse)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ReaderService/ReadQuorumCallCombo", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalReadResponse{n.id, reply, err}
}

// ReadEmpty and other methods for testing imported protos
func (c *Configuration) ReadEmpty(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (resp *ReadResponse, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadEmpty")
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

	expected := c.n
	replyChan := make(chan internalReadResponse, expected)
	for _, n := range c.nodes {
		go n.ReadEmpty(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*ReadResponse, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadEmptyQF(replyValues); quorum {
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

func (n *Node) ReadEmpty(ctx context.Context, in *empty.Empty, replyChan chan<- internalReadResponse) {
	reply := new(ReadResponse)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ReaderService/ReadEmpty", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalReadResponse{n.id, reply, err}
}

// ReadEmpty2 is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) ReadEmpty2(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (resp *empty.Empty, err error) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadEmpty2")
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

	expected := c.n
	replyChan := make(chan internalEmpty, expected)
	for _, n := range c.nodes {
		go n.ReadEmpty2(ctx, in, replyChan)
	}

	var (
		replyValues = make([]*empty.Empty, 0, expected)
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}

			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}

			replyValues = append(replyValues, r.reply)
			if resp, quorum = c.qspec.ReadEmpty2QF(replyValues); quorum {
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

func (n *Node) ReadEmpty2(ctx context.Context, in *ReadRequest, replyChan chan<- internalEmpty) {
	reply := new(empty.Empty)
	start := time.Now()
	err := n.conn.Invoke(ctx, "/dev.ReaderService/ReadEmpty2", in, reply)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		n.setLatency(time.Since(start))
	} else {
		n.setLastErr(err)
	}
	replyChan <- internalEmpty{n.id, reply, err}
}
