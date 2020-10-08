// Code generated by protoc-gen-gorums. DO NOT EDIT.

package dev

import (
	context "context"
	gorums "github.com/relab/gorums"
)

// GRPCCall plain gRPC call; testing that Gorums can ignore these, but that
// they are added to the _grpc.pb.go generated file.
func (n *Node) GRPCCall(ctx context.Context, in *Request) (resp *Response, err error) {

	cd := gorums.CallData{
		Manager:  n.mgr.Manager,
		Node:     n.Node,
		Message:  in,
		MethodID: gRPCCallMethodID,
	}

	res, err := gorums.RPCCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Response), err
}
