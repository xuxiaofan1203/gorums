// Code generated by protoc-gen-gorums. DO NOT EDIT.
// Source file to edit is: dev/storage.proto
// Template file to edit is: calltype_correctable.tmpl

package dev

import (
	"context"
	"time"

	"golang.org/x/net/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* Exported correctable method ReadCorrectable */

// ReadCorrectable asynchronously invokes a
// correctable ReadCorrectable quorum call on configuration c and returns a
// CorrectableState which can be used to inspect any replies or errors
// when available.
func (c *Configuration) ReadCorrectable(ctx context.Context, arg *ReadRequest) *CorrectableState {
	corr := &CorrectableState{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.readCorrectable(ctx, arg, corr)
	return corr
}

// Get returns the reply, level and any error associated with the
// ReadCorrectable. The method does not block until a (possibly
// itermidiate) reply or error is available. Level is set to LevelNotSet if no
// reply has yet been received. The Done or Watch methods should be used to
// ensure that a reply is available.
func (c *CorrectableState) Get() (*State, int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.State, c.level, c.err
}

// Done returns a channel that's closed when the correctable ReadCorrectable
// quorum call is done. A call is considered done when the quorum function has
// signaled that a quorum of replies was received or that the call returned an
// error.
func (c *CorrectableState) Done() <-chan struct{} {
	return c.donech
}

// Watch returns a channel that's closed when a reply or error at or above the
// specified level is available. If the call is done, the channel is closed
// disregardless of the specified level.
func (c *CorrectableState) Watch(level int) <-chan struct{} {
	ch := make(chan struct{})
	c.mu.Lock()
	if level < c.level {
		close(ch)
		c.mu.Unlock()
		return ch
	}
	c.watchers = append(c.watchers, &struct {
		level int
		ch    chan struct{}
	}{level, ch})
	c.mu.Unlock()
	return ch
}

func (c *CorrectableState) set(reply *State, level int, err error, done bool) {
	c.mu.Lock()
	if c.done {
		c.mu.Unlock()
		panic("set(...) called on a done correctable")
	}
	c.State, c.level, c.err, c.done = reply, level, err, done
	if done {
		close(c.donech)
		for _, watcher := range c.watchers {
			if watcher != nil {
				close(watcher.ch)
			}
		}
		c.mu.Unlock()
		return
	}
	for i := range c.watchers {
		if c.watchers[i] != nil && c.watchers[i].level <= level {
			close(c.watchers[i].ch)
			c.watchers[i] = nil
		}
	}
	c.mu.Unlock()
}

/* Unexported correctable method ReadCorrectable */

func (c *Configuration) readCorrectable(ctx context.Context, a *ReadRequest, resp *CorrectableState) {
	var ti traceInfo
	if c.mgr.opts.trace {
		ti.Trace = trace.New("gorums."+c.tstring()+".Sent", "ReadCorrectable")
		defer ti.Finish()

		ti.firstLine.cid = c.id
		if deadline, ok := ctx.Deadline(); ok {
			ti.firstLine.deadline = time.Until(deadline)
		}
		ti.LazyLog(&ti.firstLine, false)
		ti.LazyLog(&payload{sent: true, msg: a}, false)

		defer func() {
			ti.LazyLog(&qcresult{
				ids:   resp.NodeIDs,
				reply: resp.State,
				err:   resp.err,
			}, false)
			if resp.err != nil {
				ti.SetError()
			}
		}()
	}

	expected := c.n
	replyChan := make(chan internalState, expected)
	for _, n := range c.nodes {
		go callGRPCReadCorrectable(ctx, n, a, replyChan)
	}

	var (
		replyValues = make([]*State, 0, c.n)
		clevel      = LevelNotSet
		reply       *State
		rlevel      int
		errs        []GRPCError
		quorum      bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = append(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			if c.mgr.opts.trace {
				ti.LazyLog(&payload{sent: false, id: r.nid, msg: r.reply}, false)
			}
			replyValues = append(replyValues, r.reply)
			reply, rlevel, quorum = c.qspec.ReadCorrectableQF(replyValues)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replyValues), errs}, true)
			return
		}

		if len(errs)+len(replyValues) == expected {
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replyValues), errs}, true)
			return
		}
	}
}

func callGRPCReadCorrectable(ctx context.Context, node *Node, arg *ReadRequest, replyChan chan<- internalState) {
	reply := new(State)
	start := time.Now()
	err := grpc.Invoke(
		ctx,
		"/dev.Storage/ReadCorrectable",
		arg,
		reply,
		node.conn,
	)
	s, ok := status.FromError(err)
	if ok && (s.Code() == codes.OK || s.Code() == codes.Canceled) {
		node.setLatency(time.Since(start))
	} else {
		node.setLastErr(err)
	}
	replyChan <- internalState{node.id, reply, err}
}
