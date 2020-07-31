package gengorums

var correctableCallVariables = `
{{$context := use "context.Context" .GenFile}}
{{$opts := use "grpc.CallOption" .GenFile}}
{{$correctableOut := outType .Method $customOut}}
{{$appendFn := printf "append"}}
{{if correctableStream .Method}}
{{$appendFn = printf "appendIfNotPresent"}}
{{end}}
`

var correctableCallComment = `
{{$comments := .Method.Comments.Leading}}
{{if ne $comments ""}}
{{$comments -}}
{{else}}
{{if hasPerNodeArg .Method}}
// {{$method}} asynchronously invokes a correctable quorum call on each node
// in configuration c, with the argument returned by the provided function f
// and returns a {{$correctableOut}}, which can be used to inspect
// the reply and error when available.
// The provide per node function f takes the provided {{$in}} argument
// and returns an {{$out}} object to be passed to the given nodeID.
// The per node function f should be thread-safe.
{{else}}
// {{$method}} asynchronously invokes a correctable quorum call on each node
// in configuration c and returns a {{$correctableOut}}, which can be used
// to inspect any replies or errors when available.
{{if correctableStream .Method -}}
// This method supports server-side preliminary replies (correctable stream).
{{end -}}
{{end -}}
{{end -}}
`

var correctableCallSignature = `func (c *Configuration) {{$method}}(` +
	`ctx {{$context}}, in *{{$in}}` +
	`{{perNodeFnType .GenFile .Method ", f"}}` +
	`, opts ...{{$opts}}) ` +
	`*{{$correctableOut}} {`

var correctableCallBody = `
	corr := &{{$correctableOut}}{
		level:   LevelNotSet,
		NodeIDs: make([]uint32, 0, c.n),
		donech:  make(chan struct{}),
	}
	go c.{{unexport .Method.GoName}}(ctx, in{{perNodeArg .Method ", f"}}, corr, opts...)
	return corr
}
`

var correctableCallUnexportedSignature = `
func (c *Configuration) {{unexport .Method.GoName}}(` +
	`ctx {{$context}}, in *{{$in}}` +
	`{{perNodeFnType .GenFile .Method ", f"}}` +
	`, resp *{{$correctableOut}}, opts ...{{$opts}}) {
`

var correctableCallReply = `
	var (
		replies = make(map[uint32]*{{$out}}, c.n*2)
		clevel		= LevelNotSet
		reply		*{{$customOut}}
		rlevel		int
		errs		[]GRPCError
		quorum		bool
	)

	for {
		select {
		case r := <-replyChan:
			resp.NodeIDs = {{$appendFn}}(resp.NodeIDs, r.nid)
			if r.err != nil {
				errs = append(errs, GRPCError{r.nid, r.err})
				break
			}
			{{template "traceLazyLog"}}
			replies[r.nid] = r.reply
			reply, rlevel, quorum = c.qspec.{{$method}}QF(in, replies)
			if quorum {
				resp.set(reply, rlevel, nil, true)
				return
			}
			if rlevel > clevel {
				clevel = rlevel
				resp.set(reply, rlevel, nil, false)
			}
		case <-ctx.Done():
			resp.set(reply, clevel, QuorumCallError{ctx.Err().Error(), len(replies), errs}, true)
			return
		}
		{{- if correctableStream .Method}}
		if len(errs) == expected { // Can't rely on reply count.
		{{else}}
		if len(errs)+len(replies) == expected {
		{{end -}}
			resp.set(reply, clevel, QuorumCallError{"incomplete call", len(replies), errs}, true)
			return
		}
	}
}
`

var correctableStreamCallGrpc = `
func (n *Node) {{$method}}(ctx {{$context}}, in *{{$in}}, replyChan chan<- {{$intOut}}) {
	x := New{{serviceName .Method}}Client(n.conn)
	y, err := x.{{$method}}(ctx, in)
	if err != nil {
		replyChan <- {{$intOut}}{n.id, nil, err}
		return
	}

	for {
		reply, err := y.Recv()
		if err == {{use "io.EOF" .GenFile}} {
			return
		}
		replyChan <- {{$intOut}}{n.id, reply, err}
		if err != nil {
			return
		}
	}
}
`

var correctableVar = `
{{$correctableOut := outType .Method $customOut}}
{{$protoMessage := use "protoreflect.ProtoMessage" .GenFile}}
{{$callData := use "gorums.CorrectableCallData" .GenFile}}
{{$genFile := .GenFile}}
{{$unexportMethod := unexport .Method.GoName}}
{{$context := use "context.Context" .GenFile}}
`

var correctableSignature = `func (c *Configuration) {{$method}}(` +
	`ctx {{$context}}, in *{{$in}}` +
	`{{perNodeFnType .GenFile .Method ", f"}}) ` +
	`*{{$correctableOut}} {`

var correctableBody = `
	cd := {{$callData}}{
		Manager:  c.mgr.Manager,
		Nodes:    c.nodes,
		Message:  in,
		MethodID: {{$unexportMethod}}MethodID,
	{{- if correctableStream .Method}}
		ServerStream: true,
	{{- else}}
		ServerStream: false,
	{{- end}}
	}
	cd.QuorumFunction = func(req {{$protoMessage}}, replies map[uint32]{{$protoMessage}}) ({{$protoMessage}}, int, bool) {
		r := make(map[uint32]*{{$out}}, len(replies))
		for k, v := range replies {
			r[k] = v.(*{{$out}})
		}
		return c.qspec.{{$method}}QF(req.(*{{$in}}), r)
	}
{{- if hasPerNodeArg .Method}}
	cd.PerNodeArgFn = func(req {{$protoMessage}}, nid uint32) {{$protoMessage}} {
		return f(req.(*{{$in}}), nid)
	}
{{- end}}

	corr := {{use "gorums.CorrectableCall" $genFile}}(ctx, cd)
	return &{{$correctableOut}}{corr}
}
`

var correctableCall = commonVariables +
	correctableVar +
	correctableCallComment +
	correctableSignature +
	correctableBody
