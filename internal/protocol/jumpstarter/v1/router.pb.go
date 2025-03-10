// Copyright 2024 The Jumpstarter Authors

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: jumpstarter/v1/router.proto

package jumpstarterv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FrameType int32

const (
	FrameType_FRAME_TYPE_DATA       FrameType = 0
	FrameType_FRAME_TYPE_RST_STREAM FrameType = 3
	FrameType_FRAME_TYPE_PING       FrameType = 6
	FrameType_FRAME_TYPE_GOAWAY     FrameType = 7
)

// Enum value maps for FrameType.
var (
	FrameType_name = map[int32]string{
		0: "FRAME_TYPE_DATA",
		3: "FRAME_TYPE_RST_STREAM",
		6: "FRAME_TYPE_PING",
		7: "FRAME_TYPE_GOAWAY",
	}
	FrameType_value = map[string]int32{
		"FRAME_TYPE_DATA":       0,
		"FRAME_TYPE_RST_STREAM": 3,
		"FRAME_TYPE_PING":       6,
		"FRAME_TYPE_GOAWAY":     7,
	}
)

func (x FrameType) Enum() *FrameType {
	p := new(FrameType)
	*p = x
	return p
}

func (x FrameType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FrameType) Descriptor() protoreflect.EnumDescriptor {
	return file_jumpstarter_v1_router_proto_enumTypes[0].Descriptor()
}

func (FrameType) Type() protoreflect.EnumType {
	return &file_jumpstarter_v1_router_proto_enumTypes[0]
}

func (x FrameType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FrameType.Descriptor instead.
func (FrameType) EnumDescriptor() ([]byte, []int) {
	return file_jumpstarter_v1_router_proto_rawDescGZIP(), []int{0}
}

type StreamRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Payload       []byte                 `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	FrameType     FrameType              `protobuf:"varint,2,opt,name=frame_type,json=frameType,proto3,enum=jumpstarter.v1.FrameType" json:"frame_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	mi := &file_jumpstarter_v1_router_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jumpstarter_v1_router_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_jumpstarter_v1_router_proto_rawDescGZIP(), []int{0}
}

func (x *StreamRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *StreamRequest) GetFrameType() FrameType {
	if x != nil {
		return x.FrameType
	}
	return FrameType_FRAME_TYPE_DATA
}

type StreamResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Payload       []byte                 `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	FrameType     FrameType              `protobuf:"varint,2,opt,name=frame_type,json=frameType,proto3,enum=jumpstarter.v1.FrameType" json:"frame_type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamResponse) Reset() {
	*x = StreamResponse{}
	mi := &file_jumpstarter_v1_router_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResponse) ProtoMessage() {}

func (x *StreamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jumpstarter_v1_router_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResponse.ProtoReflect.Descriptor instead.
func (*StreamResponse) Descriptor() ([]byte, []int) {
	return file_jumpstarter_v1_router_proto_rawDescGZIP(), []int{1}
}

func (x *StreamResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *StreamResponse) GetFrameType() FrameType {
	if x != nil {
		return x.FrameType
	}
	return FrameType_FRAME_TYPE_DATA
}

var File_jumpstarter_v1_router_proto protoreflect.FileDescriptor

var file_jumpstarter_v1_router_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31,
	0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6a,
	0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x63, 0x0a,
	0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x38, 0x0a, 0x0a, 0x66, 0x72, 0x61, 0x6d,
	0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x6a,
	0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x72,
	0x61, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x22, 0x64, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x38,
	0x0a, 0x0a, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x19, 0x2e, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x66,
	0x72, 0x61, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65, 0x2a, 0x67, 0x0a, 0x09, 0x46, 0x72, 0x61, 0x6d,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x10, 0x00, 0x12, 0x19, 0x0a, 0x15, 0x46, 0x52,
	0x41, 0x4d, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x53, 0x54, 0x5f, 0x53, 0x54, 0x52,
	0x45, 0x41, 0x4d, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x50, 0x49, 0x4e, 0x47, 0x10, 0x06, 0x12, 0x15, 0x0a, 0x11, 0x46, 0x52,
	0x41, 0x4d, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x4f, 0x41, 0x57, 0x41, 0x59, 0x10,
	0x07, 0x32, 0x5c, 0x0a, 0x0d, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x4b, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x1d, 0x2e, 0x6a,
	0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6a, 0x75,
	0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42,
	0xca, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2d, 0x64, 0x65,
	0x76, 0x2f, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2d, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x4a, 0x58, 0x58, 0xaa, 0x02, 0x0e, 0x4a, 0x75,
	0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0e, 0x4a,
	0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1a,
	0x4a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47,
	0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0f, 0x4a, 0x75, 0x6d,
	0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_jumpstarter_v1_router_proto_rawDescOnce sync.Once
	file_jumpstarter_v1_router_proto_rawDescData []byte
)

func file_jumpstarter_v1_router_proto_rawDescGZIP() []byte {
	file_jumpstarter_v1_router_proto_rawDescOnce.Do(func() {
		file_jumpstarter_v1_router_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_jumpstarter_v1_router_proto_rawDesc), len(file_jumpstarter_v1_router_proto_rawDesc)))
	})
	return file_jumpstarter_v1_router_proto_rawDescData
}

var file_jumpstarter_v1_router_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_jumpstarter_v1_router_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_jumpstarter_v1_router_proto_goTypes = []any{
	(FrameType)(0),         // 0: jumpstarter.v1.FrameType
	(*StreamRequest)(nil),  // 1: jumpstarter.v1.StreamRequest
	(*StreamResponse)(nil), // 2: jumpstarter.v1.StreamResponse
}
var file_jumpstarter_v1_router_proto_depIdxs = []int32{
	0, // 0: jumpstarter.v1.StreamRequest.frame_type:type_name -> jumpstarter.v1.FrameType
	0, // 1: jumpstarter.v1.StreamResponse.frame_type:type_name -> jumpstarter.v1.FrameType
	1, // 2: jumpstarter.v1.RouterService.Stream:input_type -> jumpstarter.v1.StreamRequest
	2, // 3: jumpstarter.v1.RouterService.Stream:output_type -> jumpstarter.v1.StreamResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_jumpstarter_v1_router_proto_init() }
func file_jumpstarter_v1_router_proto_init() {
	if File_jumpstarter_v1_router_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_jumpstarter_v1_router_proto_rawDesc), len(file_jumpstarter_v1_router_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_jumpstarter_v1_router_proto_goTypes,
		DependencyIndexes: file_jumpstarter_v1_router_proto_depIdxs,
		EnumInfos:         file_jumpstarter_v1_router_proto_enumTypes,
		MessageInfos:      file_jumpstarter_v1_router_proto_msgTypes,
	}.Build()
	File_jumpstarter_v1_router_proto = out.File
	file_jumpstarter_v1_router_proto_goTypes = nil
	file_jumpstarter_v1_router_proto_depIdxs = nil
}
