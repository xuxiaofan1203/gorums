// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.16.0
// source: metadata/metadata.proto

package metadata

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NodeID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID uint32 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *NodeID) Reset() {
	*x = NodeID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metadata_metadata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeID) ProtoMessage() {}

func (x *NodeID) ProtoReflect() protoreflect.Message {
	mi := &file_metadata_metadata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeID.ProtoReflect.Descriptor instead.
func (*NodeID) Descriptor() ([]byte, []int) {
	return file_metadata_metadata_proto_rawDescGZIP(), []int{0}
}

func (x *NodeID) GetID() uint32 {
	if x != nil {
		return x.ID
	}
	return 0
}

type IPAddr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr string `protobuf:"bytes,1,opt,name=Addr,proto3" json:"Addr,omitempty"`
}

func (x *IPAddr) Reset() {
	*x = IPAddr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metadata_metadata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IPAddr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IPAddr) ProtoMessage() {}

func (x *IPAddr) ProtoReflect() protoreflect.Message {
	mi := &file_metadata_metadata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IPAddr.ProtoReflect.Descriptor instead.
func (*IPAddr) Descriptor() ([]byte, []int) {
	return file_metadata_metadata_proto_rawDescGZIP(), []int{1}
}

func (x *IPAddr) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

var File_metadata_metadata_proto protoreflect.FileDescriptor

var file_metadata_metadata_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x18, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x49, 0x44, 0x22, 0x1c, 0x0a, 0x06, 0x49, 0x50,
	0x41, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x41, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x41, 0x64, 0x64, 0x72, 0x32, 0x7c, 0x0a, 0x0c, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x54, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x49, 0x44, 0x46, 0x72,
	0x6f, 0x6d, 0x4d, 0x44, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x44, 0x22, 0x00,
	0x12, 0x34, 0x0a, 0x06, 0x57, 0x68, 0x61, 0x74, 0x49, 0x50, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x10, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x49, 0x50,
	0x41, 0x64, 0x64, 0x72, 0x22, 0x00, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x62, 0x2f, 0x67, 0x6f, 0x72, 0x75, 0x6d,
	0x73, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metadata_metadata_proto_rawDescOnce sync.Once
	file_metadata_metadata_proto_rawDescData = file_metadata_metadata_proto_rawDesc
)

func file_metadata_metadata_proto_rawDescGZIP() []byte {
	file_metadata_metadata_proto_rawDescOnce.Do(func() {
		file_metadata_metadata_proto_rawDescData = protoimpl.X.CompressGZIP(file_metadata_metadata_proto_rawDescData)
	})
	return file_metadata_metadata_proto_rawDescData
}

var file_metadata_metadata_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_metadata_metadata_proto_goTypes = []interface{}{
	(*NodeID)(nil),        // 0: metadata.NodeID
	(*IPAddr)(nil),        // 1: metadata.IPAddr
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_metadata_metadata_proto_depIdxs = []int32{
	2, // 0: metadata.MetadataTest.IDFromMD:input_type -> google.protobuf.Empty
	2, // 1: metadata.MetadataTest.WhatIP:input_type -> google.protobuf.Empty
	0, // 2: metadata.MetadataTest.IDFromMD:output_type -> metadata.NodeID
	1, // 3: metadata.MetadataTest.WhatIP:output_type -> metadata.IPAddr
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_metadata_metadata_proto_init() }
func file_metadata_metadata_proto_init() {
	if File_metadata_metadata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metadata_metadata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_metadata_metadata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IPAddr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_metadata_metadata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metadata_metadata_proto_goTypes,
		DependencyIndexes: file_metadata_metadata_proto_depIdxs,
		MessageInfos:      file_metadata_metadata_proto_msgTypes,
	}.Build()
	File_metadata_metadata_proto = out.File
	file_metadata_metadata_proto_rawDesc = nil
	file_metadata_metadata_proto_goTypes = nil
	file_metadata_metadata_proto_depIdxs = nil
}
