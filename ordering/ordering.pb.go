// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: ordering/ordering.proto

package ordering

import (
	status "google.golang.org/genproto/googleapis/rpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Metadata is sent together with application-specific message types,
// and contains information necessary for Gorums to handle the messages.
type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageID uint64         `protobuf:"varint,1,opt,name=MessageID,proto3" json:"MessageID,omitempty"`
	Method    string         `protobuf:"bytes,2,opt,name=Method,proto3" json:"Method,omitempty"`
	Status    *status.Status `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ordering_ordering_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_ordering_ordering_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_ordering_ordering_proto_rawDescGZIP(), []int{0}
}

func (x *Metadata) GetMessageID() uint64 {
	if x != nil {
		return x.MessageID
	}
	return 0
}

func (x *Metadata) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Metadata) GetStatus() *status.Status {
	if x != nil {
		return x.Status
	}
	return nil
}

var File_ordering_ordering_proto protoreflect.FileDescriptor

var file_ordering_ordering_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x1a, 0x17, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x72, 0x70, 0x63, 0x2f,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6c, 0x0a, 0x08,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x2a,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x42, 0x0a, 0x06, 0x47, 0x6f,
	0x72, 0x75, 0x6d, 0x73, 0x12, 0x38, 0x0a, 0x0a, 0x4e, 0x6f, 0x64, 0x65, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x12, 0x12, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x12, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x69, 0x6e,
	0x67, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x28, 0x01, 0x30, 0x01, 0x42, 0x22,
	0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x65, 0x6c,
	0x61, 0x62, 0x2f, 0x67, 0x6f, 0x72, 0x75, 0x6d, 0x73, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ordering_ordering_proto_rawDescOnce sync.Once
	file_ordering_ordering_proto_rawDescData = file_ordering_ordering_proto_rawDesc
)

func file_ordering_ordering_proto_rawDescGZIP() []byte {
	file_ordering_ordering_proto_rawDescOnce.Do(func() {
		file_ordering_ordering_proto_rawDescData = protoimpl.X.CompressGZIP(file_ordering_ordering_proto_rawDescData)
	})
	return file_ordering_ordering_proto_rawDescData
}

var file_ordering_ordering_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ordering_ordering_proto_goTypes = []interface{}{
	(*Metadata)(nil),      // 0: ordering.Metadata
	(*status.Status)(nil), // 1: google.rpc.Status
}
var file_ordering_ordering_proto_depIdxs = []int32{
	1, // 0: ordering.Metadata.Status:type_name -> google.rpc.Status
	0, // 1: ordering.Gorums.NodeStream:input_type -> ordering.Metadata
	0, // 2: ordering.Gorums.NodeStream:output_type -> ordering.Metadata
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ordering_ordering_proto_init() }
func file_ordering_ordering_proto_init() {
	if File_ordering_ordering_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ordering_ordering_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metadata); i {
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
			RawDescriptor: file_ordering_ordering_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ordering_ordering_proto_goTypes,
		DependencyIndexes: file_ordering_ordering_proto_depIdxs,
		MessageInfos:      file_ordering_ordering_proto_msgTypes,
	}.Build()
	File_ordering_ordering_proto = out.File
	file_ordering_ordering_proto_rawDesc = nil
	file_ordering_ordering_proto_goTypes = nil
	file_ordering_ordering_proto_depIdxs = nil
}
