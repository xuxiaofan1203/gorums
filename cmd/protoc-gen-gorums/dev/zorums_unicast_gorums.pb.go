// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.3.0-devel
// 	protoc            v3.15.3
// source: zorums.proto

package dev

import (
	context "context"
	gorums "github.com/relab/gorums"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// Unicast is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Unicast(ctx context.Context, in *Request, opts ...gorums.CallOption) {
	cd := gorums.CallData{
		Message: in,
		Method:  "dev.ZorumsService.Unicast",
	}

	n.Node.Unicast(ctx, cd, opts...)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ emptypb.Empty

// Unicast2 is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Unicast2(ctx context.Context, in *Request, opts ...gorums.CallOption) {
	cd := gorums.CallData{
		Message: in,
		Method:  "dev.ZorumsService.Unicast2",
	}

	n.Node.Unicast(ctx, cd, opts...)
}
