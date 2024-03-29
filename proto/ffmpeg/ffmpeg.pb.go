// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: ffmpeg.proto

package ffmpeg

import (
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

// Result message used to send error if occurred.
type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"` // Error message if any occurred, empty if no error.
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{0}
}

func (x *Result) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// File represents a stream of bytes which would be the equivalent of an io.ReadCloser in protobuf.
type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"` // The contents of the file as a byte array since streaming isn't natively supported.
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{1}
}

func (x *File) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

// Job message struct to be sent over the network.
type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Blender string   `protobuf:"bytes,1,opt,name=blender,proto3" json:"blender,omitempty"` // Equivalent of Blender filepath.
	Out     []string `protobuf:"bytes,2,rep,name=out,proto3" json:"out,omitempty"`         // List of file extensions.
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{2}
}

func (x *Job) GetBlender() string {
	if x != nil {
		return x.Blender
	}
	return ""
}

func (x *Job) GetOut() []string {
	if x != nil {
		return x.Out
	}
	return nil
}

type ProcessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*ProcessRequest_Args
	//	*ProcessRequest_Meta
	//	*ProcessRequest_Blob
	Data isProcessRequest_Data `protobuf_oneof:"data"`
}

func (x *ProcessRequest) Reset() {
	*x = ProcessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRequest) ProtoMessage() {}

func (x *ProcessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRequest.ProtoReflect.Descriptor instead.
func (*ProcessRequest) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{3}
}

func (m *ProcessRequest) GetData() isProcessRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ProcessRequest) GetArgs() *ProcessArgs {
	if x, ok := x.GetData().(*ProcessRequest_Args); ok {
		return x.Args
	}
	return nil
}

func (x *ProcessRequest) GetMeta() *TargetMeta {
	if x, ok := x.GetData().(*ProcessRequest_Meta); ok {
		return x.Meta
	}
	return nil
}

func (x *ProcessRequest) GetBlob() *TargetBlob {
	if x, ok := x.GetData().(*ProcessRequest_Blob); ok {
		return x.Blob
	}
	return nil
}

type isProcessRequest_Data interface {
	isProcessRequest_Data()
}

type ProcessRequest_Args struct {
	Args *ProcessArgs `protobuf:"bytes,1,opt,name=args,proto3,oneof"`
}

type ProcessRequest_Meta struct {
	Meta *TargetMeta `protobuf:"bytes,2,opt,name=meta,proto3,oneof"`
}

type ProcessRequest_Blob struct {
	Blob *TargetBlob `protobuf:"bytes,3,opt,name=blob,proto3,oneof"`
}

func (*ProcessRequest_Args) isProcessRequest_Data() {}

func (*ProcessRequest_Meta) isProcessRequest_Data() {}

func (*ProcessRequest_Blob) isProcessRequest_Data() {}

// for processing what video formats to convert to
type ProcessArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExtensionList []string `protobuf:"bytes,1,rep,name=extension_list,json=extensionList,proto3" json:"extension_list,omitempty"`
}

func (x *ProcessArgs) Reset() {
	*x = ProcessArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessArgs) ProtoMessage() {}

func (x *ProcessArgs) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessArgs.ProtoReflect.Descriptor instead.
func (*ProcessArgs) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{4}
}

func (x *ProcessArgs) GetExtensionList() []string {
	if x != nil {
		return x.ExtensionList
	}
	return nil
}

// for readings the size of the file being uploaded
type TargetMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Size      uint64 `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	Extension string `protobuf:"bytes,2,opt,name=extension,proto3" json:"extension,omitempty"`
}

func (x *TargetMeta) Reset() {
	*x = TargetMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TargetMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TargetMeta) ProtoMessage() {}

func (x *TargetMeta) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TargetMeta.ProtoReflect.Descriptor instead.
func (*TargetMeta) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{5}
}

func (x *TargetMeta) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *TargetMeta) GetExtension() string {
	if x != nil {
		return x.Extension
	}
	return ""
}

// for uploading a file
type TargetBlob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data      []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Extension string `protobuf:"bytes,2,opt,name=extension,proto3" json:"extension,omitempty"`
}

func (x *TargetBlob) Reset() {
	*x = TargetBlob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TargetBlob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TargetBlob) ProtoMessage() {}

func (x *TargetBlob) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TargetBlob.ProtoReflect.Descriptor instead.
func (*TargetBlob) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{6}
}

func (x *TargetBlob) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *TargetBlob) GetExtension() string {
	if x != nil {
		return x.Extension
	}
	return ""
}

type Log struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Log string `protobuf:"bytes,1,opt,name=log,proto3" json:"log,omitempty"`
}

func (x *Log) Reset() {
	*x = Log{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Log) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Log) ProtoMessage() {}

func (x *Log) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Log.ProtoReflect.Descriptor instead.
func (*Log) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{7}
}

func (x *Log) GetLog() string {
	if x != nil {
		return x.Log
	}
	return ""
}

// for downloading files by extension
type ProcessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*ProcessResponse_Blob
	//	*ProcessResponse_Log
	Data isProcessResponse_Data `protobuf_oneof:"data"`
}

func (x *ProcessResponse) Reset() {
	*x = ProcessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ffmpeg_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessResponse) ProtoMessage() {}

func (x *ProcessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ffmpeg_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessResponse.ProtoReflect.Descriptor instead.
func (*ProcessResponse) Descriptor() ([]byte, []int) {
	return file_ffmpeg_proto_rawDescGZIP(), []int{8}
}

func (m *ProcessResponse) GetData() isProcessResponse_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ProcessResponse) GetBlob() *TargetBlob {
	if x, ok := x.GetData().(*ProcessResponse_Blob); ok {
		return x.Blob
	}
	return nil
}

func (x *ProcessResponse) GetLog() *Log {
	if x, ok := x.GetData().(*ProcessResponse_Log); ok {
		return x.Log
	}
	return nil
}

type isProcessResponse_Data interface {
	isProcessResponse_Data()
}

type ProcessResponse_Blob struct {
	Blob *TargetBlob `protobuf:"bytes,1,opt,name=blob,proto3,oneof"`
}

type ProcessResponse_Log struct {
	Log *Log `protobuf:"bytes,2,opt,name=log,proto3,oneof"`
}

func (*ProcessResponse_Blob) isProcessResponse_Data() {}

func (*ProcessResponse_Log) isProcessResponse_Data() {}

var File_ffmpeg_proto protoreflect.FileDescriptor

var file_ffmpeg_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x22, 0x1e, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x20, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x31, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12,
	0x18, 0x0a, 0x07, 0x62, 0x6c, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x62, 0x6c, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x6f, 0x75, 0x74,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x6f, 0x75, 0x74, 0x22, 0x97, 0x01, 0x0a, 0x0e,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29,
	0x0a, 0x04, 0x61, 0x72, 0x67, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x66,
	0x66, 0x6d, 0x70, 0x65, 0x67, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x41, 0x72, 0x67,
	0x73, 0x48, 0x00, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x12, 0x28, 0x0a, 0x04, 0x6d, 0x65, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x66, 0x6d, 0x70, 0x65, 0x67,
	0x2e, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x48, 0x00, 0x52, 0x04, 0x6d,
	0x65, 0x74, 0x61, 0x12, 0x28, 0x0a, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x2e, 0x54, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x42, 0x6c, 0x6f, 0x62, 0x48, 0x00, 0x52, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x42, 0x06, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x34, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x41, 0x72, 0x67, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x3e, 0x0a, 0x0a, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3e, 0x0a, 0x0a, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x62, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a,
	0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x17, 0x0a, 0x03, 0x4c,
	0x6f, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6c, 0x6f, 0x67, 0x22, 0x64, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x62, 0x6c, 0x6f, 0x62, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x2e, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x62, 0x48, 0x00, 0x52, 0x04, 0x62, 0x6c, 0x6f,
	0x62, 0x12, 0x1f, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b,
	0x2e, 0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x2e, 0x4c, 0x6f, 0x67, 0x48, 0x00, 0x52, 0x03, 0x6c,
	0x6f, 0x67, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x4c, 0x0a, 0x0a, 0x4a, 0x6f,
	0x62, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x3e, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x16, 0x2e, 0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x2e, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x66, 0x66,
	0x6d, 0x70, 0x65, 0x67, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x70, 0x61, 0x64, 0x2f,
	0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x2d, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x66, 0x6d, 0x70, 0x65, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_ffmpeg_proto_rawDescOnce sync.Once
	file_ffmpeg_proto_rawDescData = file_ffmpeg_proto_rawDesc
)

func file_ffmpeg_proto_rawDescGZIP() []byte {
	file_ffmpeg_proto_rawDescOnce.Do(func() {
		file_ffmpeg_proto_rawDescData = protoimpl.X.CompressGZIP(file_ffmpeg_proto_rawDescData)
	})
	return file_ffmpeg_proto_rawDescData
}

var file_ffmpeg_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_ffmpeg_proto_goTypes = []interface{}{
	(*Result)(nil),          // 0: ffmpeg.Result
	(*File)(nil),            // 1: ffmpeg.File
	(*Job)(nil),             // 2: ffmpeg.Job
	(*ProcessRequest)(nil),  // 3: ffmpeg.ProcessRequest
	(*ProcessArgs)(nil),     // 4: ffmpeg.ProcessArgs
	(*TargetMeta)(nil),      // 5: ffmpeg.TargetMeta
	(*TargetBlob)(nil),      // 6: ffmpeg.TargetBlob
	(*Log)(nil),             // 7: ffmpeg.Log
	(*ProcessResponse)(nil), // 8: ffmpeg.ProcessResponse
}
var file_ffmpeg_proto_depIdxs = []int32{
	4, // 0: ffmpeg.ProcessRequest.args:type_name -> ffmpeg.ProcessArgs
	5, // 1: ffmpeg.ProcessRequest.meta:type_name -> ffmpeg.TargetMeta
	6, // 2: ffmpeg.ProcessRequest.blob:type_name -> ffmpeg.TargetBlob
	6, // 3: ffmpeg.ProcessResponse.blob:type_name -> ffmpeg.TargetBlob
	7, // 4: ffmpeg.ProcessResponse.log:type_name -> ffmpeg.Log
	3, // 5: ffmpeg.JobManager.Process:input_type -> ffmpeg.ProcessRequest
	8, // 6: ffmpeg.JobManager.Process:output_type -> ffmpeg.ProcessResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_ffmpeg_proto_init() }
func file_ffmpeg_proto_init() {
	if File_ffmpeg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ffmpeg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
		file_ffmpeg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
		file_ffmpeg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
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
		file_ffmpeg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRequest); i {
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
		file_ffmpeg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessArgs); i {
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
		file_ffmpeg_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TargetMeta); i {
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
		file_ffmpeg_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TargetBlob); i {
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
		file_ffmpeg_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Log); i {
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
		file_ffmpeg_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessResponse); i {
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
	file_ffmpeg_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*ProcessRequest_Args)(nil),
		(*ProcessRequest_Meta)(nil),
		(*ProcessRequest_Blob)(nil),
	}
	file_ffmpeg_proto_msgTypes[8].OneofWrappers = []interface{}{
		(*ProcessResponse_Blob)(nil),
		(*ProcessResponse_Log)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ffmpeg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ffmpeg_proto_goTypes,
		DependencyIndexes: file_ffmpeg_proto_depIdxs,
		MessageInfos:      file_ffmpeg_proto_msgTypes,
	}.Build()
	File_ffmpeg_proto = out.File
	file_ffmpeg_proto_rawDesc = nil
	file_ffmpeg_proto_goTypes = nil
	file_ffmpeg_proto_depIdxs = nil
}
