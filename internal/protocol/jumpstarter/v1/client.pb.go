// Copyright 2024 The Jumpstarter Authors

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: jumpstarter/v1/client.proto

package jumpstarterv1

import (
	_ "github.com/jumpstarter-dev/jumpstarter-controller/google/api"
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

type Exporter struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Labels        map[string]string      `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Exporter) Reset() {
	*x = Exporter{}
	mi := &file_jumpstarter_v1_client_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Exporter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Exporter) ProtoMessage() {}

func (x *Exporter) ProtoReflect() protoreflect.Message {
	mi := &file_jumpstarter_v1_client_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Exporter.ProtoReflect.Descriptor instead.
func (*Exporter) Descriptor() ([]byte, []int) {
	return file_jumpstarter_v1_client_proto_rawDescGZIP(), []int{0}
}

func (x *Exporter) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Exporter) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type GetExporterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetExporterRequest) Reset() {
	*x = GetExporterRequest{}
	mi := &file_jumpstarter_v1_client_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetExporterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetExporterRequest) ProtoMessage() {}

func (x *GetExporterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jumpstarter_v1_client_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetExporterRequest.ProtoReflect.Descriptor instead.
func (*GetExporterRequest) Descriptor() ([]byte, []int) {
	return file_jumpstarter_v1_client_proto_rawDescGZIP(), []int{1}
}

func (x *GetExporterRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ListExportersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Parent        string                 `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	PageSize      int32                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageToken     string                 `protobuf:"bytes,3,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	Filter        string                 `protobuf:"bytes,4,opt,name=filter,proto3" json:"filter,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListExportersRequest) Reset() {
	*x = ListExportersRequest{}
	mi := &file_jumpstarter_v1_client_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListExportersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListExportersRequest) ProtoMessage() {}

func (x *ListExportersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_jumpstarter_v1_client_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListExportersRequest.ProtoReflect.Descriptor instead.
func (*ListExportersRequest) Descriptor() ([]byte, []int) {
	return file_jumpstarter_v1_client_proto_rawDescGZIP(), []int{2}
}

func (x *ListExportersRequest) GetParent() string {
	if x != nil {
		return x.Parent
	}
	return ""
}

func (x *ListExportersRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListExportersRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListExportersRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

type ListExportersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Exporters     []*Exporter            `protobuf:"bytes,1,rep,name=exporters,proto3" json:"exporters,omitempty"`
	NextPageToken string                 `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListExportersResponse) Reset() {
	*x = ListExportersResponse{}
	mi := &file_jumpstarter_v1_client_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListExportersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListExportersResponse) ProtoMessage() {}

func (x *ListExportersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_jumpstarter_v1_client_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListExportersResponse.ProtoReflect.Descriptor instead.
func (*ListExportersResponse) Descriptor() ([]byte, []int) {
	return file_jumpstarter_v1_client_proto_rawDescGZIP(), []int{3}
}

func (x *ListExportersResponse) GetExporters() []*Exporter {
	if x != nil {
		return x.Exporters
	}
	return nil
}

func (x *ListExportersResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

var File_jumpstarter_v1_client_proto protoreflect.FileDescriptor

var file_jumpstarter_v1_client_proto_rawDesc = string([]byte{
	0x0a, 0x1b, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x6a,
	0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xfd, 0x01, 0x0a, 0x08, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x12, 0x17, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x08,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3c, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x3a,
	0x5f, 0xea, 0x41, 0x5c, 0x0a, 0x18, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65,
	0x72, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x12, 0x2b,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x2f, 0x7b, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x7d, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73,
	0x2f, 0x7b, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x7d, 0x2a, 0x09, 0x65, 0x78, 0x70,
	0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x32, 0x08, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72,
	0x22, 0x4a, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x20, 0xe0, 0x41, 0x02, 0xfa, 0x41, 0x1a, 0x0a, 0x18, 0x6a, 0x75,
	0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x45, 0x78,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xb3, 0x01, 0x0a,
	0x14, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x20, 0xe0, 0x41, 0x02, 0xfa, 0x41, 0x1a, 0x12, 0x18, 0x6a,
	0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x45,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12,
	0x20, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x03, 0xe0, 0x41, 0x01, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x22, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x01, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1b, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x01, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x22, 0x77, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x09, 0x65,
	0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52, 0x09, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x72, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65,
	0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65,
	0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xa5, 0x02, 0x0a, 0x0d,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7f, 0x0a,
	0x0b, 0x47, 0x65, 0x74, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x6a,
	0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x18, 0x2e, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x22, 0x32, 0xda, 0x41, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x12, 0x23, 0x2f, 0x76, 0x31, 0x2f, 0x7b,
	0x6e, 0x61, 0x6d, 0x65, 0x3d, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x2f,
	0x2a, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x2a, 0x7d, 0x12, 0x92,
	0x01, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73,
	0x12, 0x24, 0x2e, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x78, 0x70, 0x6f,
	0x72, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x34, 0xda,
	0x41, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x12, 0x23,
	0x2f, 0x76, 0x31, 0x2f, 0x7b, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x3d, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x2f, 0x2a, 0x7d, 0x2f, 0x65, 0x78, 0x70, 0x6f, 0x72, 0x74,
	0x65, 0x72, 0x73, 0x42, 0xca, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x6a, 0x75, 0x6d, 0x70,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65,
	0x72, 0x2d, 0x64, 0x65, 0x76, 0x2f, 0x6a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65,
	0x72, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x6a, 0x75, 0x6d,
	0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x6a, 0x75, 0x6d, 0x70,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x4a, 0x58, 0x58, 0xaa,
	0x02, 0x0e, 0x4a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0e, 0x4a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x1a, 0x4a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0f, 0x4a, 0x75, 0x6d, 0x70, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_jumpstarter_v1_client_proto_rawDescOnce sync.Once
	file_jumpstarter_v1_client_proto_rawDescData []byte
)

func file_jumpstarter_v1_client_proto_rawDescGZIP() []byte {
	file_jumpstarter_v1_client_proto_rawDescOnce.Do(func() {
		file_jumpstarter_v1_client_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_jumpstarter_v1_client_proto_rawDesc), len(file_jumpstarter_v1_client_proto_rawDesc)))
	})
	return file_jumpstarter_v1_client_proto_rawDescData
}

var file_jumpstarter_v1_client_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_jumpstarter_v1_client_proto_goTypes = []any{
	(*Exporter)(nil),              // 0: jumpstarter.v1.Exporter
	(*GetExporterRequest)(nil),    // 1: jumpstarter.v1.GetExporterRequest
	(*ListExportersRequest)(nil),  // 2: jumpstarter.v1.ListExportersRequest
	(*ListExportersResponse)(nil), // 3: jumpstarter.v1.ListExportersResponse
	nil,                           // 4: jumpstarter.v1.Exporter.LabelsEntry
}
var file_jumpstarter_v1_client_proto_depIdxs = []int32{
	4, // 0: jumpstarter.v1.Exporter.labels:type_name -> jumpstarter.v1.Exporter.LabelsEntry
	0, // 1: jumpstarter.v1.ListExportersResponse.exporters:type_name -> jumpstarter.v1.Exporter
	1, // 2: jumpstarter.v1.ClientService.GetExporter:input_type -> jumpstarter.v1.GetExporterRequest
	2, // 3: jumpstarter.v1.ClientService.ListExporters:input_type -> jumpstarter.v1.ListExportersRequest
	0, // 4: jumpstarter.v1.ClientService.GetExporter:output_type -> jumpstarter.v1.Exporter
	3, // 5: jumpstarter.v1.ClientService.ListExporters:output_type -> jumpstarter.v1.ListExportersResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_jumpstarter_v1_client_proto_init() }
func file_jumpstarter_v1_client_proto_init() {
	if File_jumpstarter_v1_client_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_jumpstarter_v1_client_proto_rawDesc), len(file_jumpstarter_v1_client_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_jumpstarter_v1_client_proto_goTypes,
		DependencyIndexes: file_jumpstarter_v1_client_proto_depIdxs,
		MessageInfos:      file_jumpstarter_v1_client_proto_msgTypes,
	}.Build()
	File_jumpstarter_v1_client_proto = out.File
	file_jumpstarter_v1_client_proto_goTypes = nil
	file_jumpstarter_v1_client_proto_depIdxs = nil
}
