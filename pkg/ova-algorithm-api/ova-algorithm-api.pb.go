// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: ova-algorithm-api/ova-algorithm-api.proto

package ova_algorithm_api

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

type AlgorithmIdV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *AlgorithmIdV1) Reset() {
	*x = AlgorithmIdV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlgorithmIdV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlgorithmIdV1) ProtoMessage() {}

func (x *AlgorithmIdV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlgorithmIdV1.ProtoReflect.Descriptor instead.
func (*AlgorithmIdV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{0}
}

func (x *AlgorithmIdV1) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type AlgorithmV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Subject     string `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *AlgorithmV1) Reset() {
	*x = AlgorithmV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlgorithmV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlgorithmV1) ProtoMessage() {}

func (x *AlgorithmV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlgorithmV1.ProtoReflect.Descriptor instead.
func (*AlgorithmV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{1}
}

func (x *AlgorithmV1) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AlgorithmV1) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *AlgorithmV1) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type AlgorithmValueV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject     string `protobuf:"bytes,2,opt,name=subject,proto3" json:"subject,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *AlgorithmValueV1) Reset() {
	*x = AlgorithmValueV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AlgorithmValueV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AlgorithmValueV1) ProtoMessage() {}

func (x *AlgorithmValueV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AlgorithmValueV1.ProtoReflect.Descriptor instead.
func (*AlgorithmValueV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{2}
}

func (x *AlgorithmValueV1) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *AlgorithmValueV1) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type CreateAlgorithmRequestV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *AlgorithmValueV1 `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *CreateAlgorithmRequestV1) Reset() {
	*x = CreateAlgorithmRequestV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAlgorithmRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAlgorithmRequestV1) ProtoMessage() {}

func (x *CreateAlgorithmRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAlgorithmRequestV1.ProtoReflect.Descriptor instead.
func (*CreateAlgorithmRequestV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{3}
}

func (x *CreateAlgorithmRequestV1) GetBody() *AlgorithmValueV1 {
	if x != nil {
		return x.Body
	}
	return nil
}

type DescribeAlgorithmRequestV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *AlgorithmIdV1 `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *DescribeAlgorithmRequestV1) Reset() {
	*x = DescribeAlgorithmRequestV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeAlgorithmRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeAlgorithmRequestV1) ProtoMessage() {}

func (x *DescribeAlgorithmRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeAlgorithmRequestV1.ProtoReflect.Descriptor instead.
func (*DescribeAlgorithmRequestV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{4}
}

func (x *DescribeAlgorithmRequestV1) GetBody() *AlgorithmIdV1 {
	if x != nil {
		return x.Body
	}
	return nil
}

type DescribeAlgorithmResponseV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *AlgorithmV1 `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *DescribeAlgorithmResponseV1) Reset() {
	*x = DescribeAlgorithmResponseV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DescribeAlgorithmResponseV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DescribeAlgorithmResponseV1) ProtoMessage() {}

func (x *DescribeAlgorithmResponseV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DescribeAlgorithmResponseV1.ProtoReflect.Descriptor instead.
func (*DescribeAlgorithmResponseV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{5}
}

func (x *DescribeAlgorithmResponseV1) GetBody() *AlgorithmV1 {
	if x != nil {
		return x.Body
	}
	return nil
}

type ListAlgorithmsRequestV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset *AlgorithmIdV1 `protobuf:"bytes,1,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit  int64          `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *ListAlgorithmsRequestV1) Reset() {
	*x = ListAlgorithmsRequestV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAlgorithmsRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAlgorithmsRequestV1) ProtoMessage() {}

func (x *ListAlgorithmsRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAlgorithmsRequestV1.ProtoReflect.Descriptor instead.
func (*ListAlgorithmsRequestV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{6}
}

func (x *ListAlgorithmsRequestV1) GetOffset() *AlgorithmIdV1 {
	if x != nil {
		return x.Offset
	}
	return nil
}

func (x *ListAlgorithmsRequestV1) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type ListAlgorithmsResponseV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body []*AlgorithmV1 `protobuf:"bytes,1,rep,name=body,proto3" json:"body,omitempty"`
}

func (x *ListAlgorithmsResponseV1) Reset() {
	*x = ListAlgorithmsResponseV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAlgorithmsResponseV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAlgorithmsResponseV1) ProtoMessage() {}

func (x *ListAlgorithmsResponseV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAlgorithmsResponseV1.ProtoReflect.Descriptor instead.
func (*ListAlgorithmsResponseV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{7}
}

func (x *ListAlgorithmsResponseV1) GetBody() []*AlgorithmV1 {
	if x != nil {
		return x.Body
	}
	return nil
}

type RemoveAlgorithmRequestV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *AlgorithmIdV1 `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *RemoveAlgorithmRequestV1) Reset() {
	*x = RemoveAlgorithmRequestV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveAlgorithmRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveAlgorithmRequestV1) ProtoMessage() {}

func (x *RemoveAlgorithmRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveAlgorithmRequestV1.ProtoReflect.Descriptor instead.
func (*RemoveAlgorithmRequestV1) Descriptor() ([]byte, []int) {
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP(), []int{8}
}

func (x *RemoveAlgorithmRequestV1) GetBody() *AlgorithmIdV1 {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_ova_algorithm_api_ova_algorithm_api_proto protoreflect.FileDescriptor

var file_ova_algorithm_api_ova_algorithm_api_proto_rawDesc = []byte{
	0x0a, 0x29, 0x6f, 0x76, 0x61, 0x2d, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2d,
	0x61, 0x70, 0x69, 0x2f, 0x6f, 0x76, 0x61, 0x2d, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x2d, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x6f, 0x76, 0x61,
	0x2e, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f, 0x0a, 0x0d, 0x41,
	0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x56, 0x31, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x59, 0x0a, 0x0b,
	0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x56, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x4e, 0x0a, 0x10, 0x41, 0x6c, 0x67, 0x6f, 0x72,
	0x69, 0x74, 0x68, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x56, 0x31, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x53, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x56, 0x31, 0x12, 0x37, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x23, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x56, 0x31, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x52, 0x0a, 0x1a,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x12, 0x34, 0x0a, 0x04, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61,
	0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x56, 0x31, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x22, 0x51, 0x0a, 0x1b, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x41, 0x6c, 0x67, 0x6f,
	0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31, 0x12,
	0x32, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x6f, 0x76, 0x61, 0x2e, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x56, 0x31, 0x52, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x22, 0x69, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x67, 0x6f, 0x72,
	0x69, 0x74, 0x68, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x12, 0x38,
	0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20,
	0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x56, 0x31,
	0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x4e,
	0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31, 0x12, 0x32, 0x0a, 0x04, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61,
	0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x56, 0x31, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x50,
	0x0a, 0x18, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x12, 0x34, 0x0a, 0x04, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61,
	0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x49, 0x64, 0x56, 0x31, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x32, 0xb0, 0x03, 0x0a, 0x0f, 0x4f, 0x76, 0x61, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68,
	0x6d, 0x41, 0x70, 0x69, 0x12, 0x5a, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x6c,
	0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x56, 0x31, 0x12, 0x2b, 0x2e, 0x6f, 0x76, 0x61, 0x2e,
	0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x76, 0x0a, 0x13, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x41, 0x6c, 0x67, 0x6f,
	0x72, 0x69, 0x74, 0x68, 0x6d, 0x56, 0x31, 0x12, 0x2d, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61, 0x6c,
	0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x62, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x1a, 0x2e, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x62, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x56, 0x31, 0x22, 0x00, 0x12, 0x6d, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74,
	0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x73, 0x56, 0x31, 0x12, 0x2a, 0x2e, 0x6f,
	0x76, 0x61, 0x2e, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x1a, 0x2b, 0x2e, 0x6f, 0x76, 0x61, 0x2e, 0x61,
	0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x56, 0x31, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x11, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x56, 0x31, 0x12, 0x2b, 0x2e, 0x6f,
	0x76, 0x61, 0x2e, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x31, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6f, 0x7a, 0x6f, 0x6e, 0x76, 0x61, 0x2f, 0x6f, 0x76, 0x61, 0x2d, 0x61, 0x6c, 0x67,
	0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6f,
	0x76, 0x61, 0x2d, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x2d, 0x61, 0x70, 0x69,
	0x3b, 0x6f, 0x76, 0x61, 0x5f, 0x61, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x5f, 0x61,
	0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ova_algorithm_api_ova_algorithm_api_proto_rawDescOnce sync.Once
	file_ova_algorithm_api_ova_algorithm_api_proto_rawDescData = file_ova_algorithm_api_ova_algorithm_api_proto_rawDesc
)

func file_ova_algorithm_api_ova_algorithm_api_proto_rawDescGZIP() []byte {
	file_ova_algorithm_api_ova_algorithm_api_proto_rawDescOnce.Do(func() {
		file_ova_algorithm_api_ova_algorithm_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_ova_algorithm_api_ova_algorithm_api_proto_rawDescData)
	})
	return file_ova_algorithm_api_ova_algorithm_api_proto_rawDescData
}

var file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_ova_algorithm_api_ova_algorithm_api_proto_goTypes = []interface{}{
	(*AlgorithmIdV1)(nil),               // 0: ova.algorithm.api.AlgorithmIdV1
	(*AlgorithmV1)(nil),                 // 1: ova.algorithm.api.AlgorithmV1
	(*AlgorithmValueV1)(nil),            // 2: ova.algorithm.api.AlgorithmValueV1
	(*CreateAlgorithmRequestV1)(nil),    // 3: ova.algorithm.api.CreateAlgorithmRequestV1
	(*DescribeAlgorithmRequestV1)(nil),  // 4: ova.algorithm.api.DescribeAlgorithmRequestV1
	(*DescribeAlgorithmResponseV1)(nil), // 5: ova.algorithm.api.DescribeAlgorithmResponseV1
	(*ListAlgorithmsRequestV1)(nil),     // 6: ova.algorithm.api.ListAlgorithmsRequestV1
	(*ListAlgorithmsResponseV1)(nil),    // 7: ova.algorithm.api.ListAlgorithmsResponseV1
	(*RemoveAlgorithmRequestV1)(nil),    // 8: ova.algorithm.api.RemoveAlgorithmRequestV1
	(*emptypb.Empty)(nil),               // 9: google.protobuf.Empty
}
var file_ova_algorithm_api_ova_algorithm_api_proto_depIdxs = []int32{
	2,  // 0: ova.algorithm.api.CreateAlgorithmRequestV1.body:type_name -> ova.algorithm.api.AlgorithmValueV1
	0,  // 1: ova.algorithm.api.DescribeAlgorithmRequestV1.body:type_name -> ova.algorithm.api.AlgorithmIdV1
	1,  // 2: ova.algorithm.api.DescribeAlgorithmResponseV1.body:type_name -> ova.algorithm.api.AlgorithmV1
	0,  // 3: ova.algorithm.api.ListAlgorithmsRequestV1.offset:type_name -> ova.algorithm.api.AlgorithmIdV1
	1,  // 4: ova.algorithm.api.ListAlgorithmsResponseV1.body:type_name -> ova.algorithm.api.AlgorithmV1
	0,  // 5: ova.algorithm.api.RemoveAlgorithmRequestV1.body:type_name -> ova.algorithm.api.AlgorithmIdV1
	3,  // 6: ova.algorithm.api.OvaAlgorithmApi.CreateAlgorithmV1:input_type -> ova.algorithm.api.CreateAlgorithmRequestV1
	4,  // 7: ova.algorithm.api.OvaAlgorithmApi.DescribeAlgorithmV1:input_type -> ova.algorithm.api.DescribeAlgorithmRequestV1
	6,  // 8: ova.algorithm.api.OvaAlgorithmApi.ListAlgorithmsV1:input_type -> ova.algorithm.api.ListAlgorithmsRequestV1
	8,  // 9: ova.algorithm.api.OvaAlgorithmApi.RemoveAlgorithmV1:input_type -> ova.algorithm.api.RemoveAlgorithmRequestV1
	9,  // 10: ova.algorithm.api.OvaAlgorithmApi.CreateAlgorithmV1:output_type -> google.protobuf.Empty
	5,  // 11: ova.algorithm.api.OvaAlgorithmApi.DescribeAlgorithmV1:output_type -> ova.algorithm.api.DescribeAlgorithmResponseV1
	7,  // 12: ova.algorithm.api.OvaAlgorithmApi.ListAlgorithmsV1:output_type -> ova.algorithm.api.ListAlgorithmsResponseV1
	9,  // 13: ova.algorithm.api.OvaAlgorithmApi.RemoveAlgorithmV1:output_type -> google.protobuf.Empty
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_ova_algorithm_api_ova_algorithm_api_proto_init() }
func file_ova_algorithm_api_ova_algorithm_api_proto_init() {
	if File_ova_algorithm_api_ova_algorithm_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlgorithmIdV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlgorithmV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AlgorithmValueV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAlgorithmRequestV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeAlgorithmRequestV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DescribeAlgorithmResponseV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAlgorithmsRequestV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAlgorithmsResponseV1); i {
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
		file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveAlgorithmRequestV1); i {
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
			RawDescriptor: file_ova_algorithm_api_ova_algorithm_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ova_algorithm_api_ova_algorithm_api_proto_goTypes,
		DependencyIndexes: file_ova_algorithm_api_ova_algorithm_api_proto_depIdxs,
		MessageInfos:      file_ova_algorithm_api_ova_algorithm_api_proto_msgTypes,
	}.Build()
	File_ova_algorithm_api_ova_algorithm_api_proto = out.File
	file_ova_algorithm_api_ova_algorithm_api_proto_rawDesc = nil
	file_ova_algorithm_api_ova_algorithm_api_proto_goTypes = nil
	file_ova_algorithm_api_ova_algorithm_api_proto_depIdxs = nil
}
