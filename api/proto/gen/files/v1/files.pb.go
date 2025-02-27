// 版本

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: api/proto/files/v1/files.proto

// 包名

package filesv1

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

type File struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                                   // 文件id
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`             // 用户id
	Url           string                 `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`                                  // 文件的minio路径
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`                                // 文件名称
	Content       string                 `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`                          // 文件内容
	CreateTime    int64                  `protobuf:"varint,6,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"` // 创建时间
	UpdateTime    int64                  `protobuf:"varint,7,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"` // 更新时间
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *File) Reset() {
	*x = File{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[0]
	if x != nil {
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
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{0}
}

func (x *File) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *File) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *File) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *File) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *File) GetUpdateTime() int64 {
	if x != nil {
		return x.UpdateTime
	}
	return 0
}

// 根据文件id，获取文件
// 请求
type GetFileByUserIdRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          *File                  `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	Limit         int64                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset        int64                  `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileByUserIdRequest) Reset() {
	*x = GetFileByUserIdRequest{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileByUserIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileByUserIdRequest) ProtoMessage() {}

func (x *GetFileByUserIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileByUserIdRequest.ProtoReflect.Descriptor instead.
func (*GetFileByUserIdRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{1}
}

func (x *GetFileByUserIdRequest) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

func (x *GetFileByUserIdRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetFileByUserIdRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

// 响应
type GetFileByUserIdResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          []*File                `protobuf:"bytes,1,rep,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileByUserIdResponse) Reset() {
	*x = GetFileByUserIdResponse{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileByUserIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileByUserIdResponse) ProtoMessage() {}

func (x *GetFileByUserIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileByUserIdResponse.ProtoReflect.Descriptor instead.
func (*GetFileByUserIdResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{2}
}

func (x *GetFileByUserIdResponse) GetFile() []*File {
	if x != nil {
		return x.File
	}
	return nil
}

// 根据文件名字，获取文件
// 请求
type GetFileByNameRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          *File                  `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileByNameRequest) Reset() {
	*x = GetFileByNameRequest{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileByNameRequest) ProtoMessage() {}

func (x *GetFileByNameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileByNameRequest.ProtoReflect.Descriptor instead.
func (*GetFileByNameRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{3}
}

func (x *GetFileByNameRequest) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

// 响应
type GetFileByNameResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          *File                  `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetFileByNameResponse) Reset() {
	*x = GetFileByNameResponse{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetFileByNameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFileByNameResponse) ProtoMessage() {}

func (x *GetFileByNameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFileByNameResponse.ProtoReflect.Descriptor instead.
func (*GetFileByNameResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{4}
}

func (x *GetFileByNameResponse) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

// 插入文件
// 请求
type CreateFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          *File                  `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateFileRequest) Reset() {
	*x = CreateFileRequest{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFileRequest) ProtoMessage() {}

func (x *CreateFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFileRequest.ProtoReflect.Descriptor instead.
func (*CreateFileRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{5}
}

func (x *CreateFileRequest) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

// 响应
type CreateFileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateFileResponse) Reset() {
	*x = CreateFileResponse{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFileResponse) ProtoMessage() {}

func (x *CreateFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFileResponse.ProtoReflect.Descriptor instead.
func (*CreateFileResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{6}
}

// 更新文件
// 请求
type UpdateFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          *File                  `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateFileRequest) Reset() {
	*x = UpdateFileRequest{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateFileRequest) ProtoMessage() {}

func (x *UpdateFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateFileRequest.ProtoReflect.Descriptor instead.
func (*UpdateFileRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateFileRequest) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

// 响应
type UpdateFileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateFileResponse) Reset() {
	*x = UpdateFileResponse{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateFileResponse) ProtoMessage() {}

func (x *UpdateFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateFileResponse.ProtoReflect.Descriptor instead.
func (*UpdateFileResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{8}
}

// 删除文件
// 请求
type DeleteFileRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	File          *File                  `protobuf:"bytes,1,opt,name=file,proto3" json:"file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileRequest) Reset() {
	*x = DeleteFileRequest{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileRequest) ProtoMessage() {}

func (x *DeleteFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileRequest.ProtoReflect.Descriptor instead.
func (*DeleteFileRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteFileRequest) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

type DeleteFileResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteFileResponse) Reset() {
	*x = DeleteFileResponse{}
	mi := &file_api_proto_files_v1_files_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileResponse) ProtoMessage() {}

func (x *DeleteFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_files_v1_files_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileResponse.ProtoReflect.Descriptor instead.
func (*DeleteFileResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_files_v1_files_proto_rawDescGZIP(), []int{10}
}

var File_api_proto_files_v1_files_proto protoreflect.FileDescriptor

var file_api_proto_files_v1_files_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65,
	0x73, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x22, 0xb1, 0x01, 0x0a, 0x04, 0x46,
	0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x6a,
	0x0a, 0x16, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x3d, 0x0a, 0x17, 0x47, 0x65,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x3a, 0x0a, 0x14, 0x47, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x3b, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22,
	0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69,
	0x6c, 0x65, 0x22, 0x37, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x37, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x37, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0x93, 0x03, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x47, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1b,
	0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1b, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x47, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x1b, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e,
	0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x56, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x20,
	0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x21, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x79,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x85, 0x01, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x46, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x14, 0x46, 0x69, 0x6c,
	0x65, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x09, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_files_v1_files_proto_rawDescOnce sync.Once
	file_api_proto_files_v1_files_proto_rawDescData = file_api_proto_files_v1_files_proto_rawDesc
)

func file_api_proto_files_v1_files_proto_rawDescGZIP() []byte {
	file_api_proto_files_v1_files_proto_rawDescOnce.Do(func() {
		file_api_proto_files_v1_files_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_files_v1_files_proto_rawDescData)
	})
	return file_api_proto_files_v1_files_proto_rawDescData
}

var file_api_proto_files_v1_files_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_proto_files_v1_files_proto_goTypes = []any{
	(*File)(nil),                    // 0: files.v1.File
	(*GetFileByUserIdRequest)(nil),  // 1: files.v1.GetFileByUserIdRequest
	(*GetFileByUserIdResponse)(nil), // 2: files.v1.GetFileByUserIdResponse
	(*GetFileByNameRequest)(nil),    // 3: files.v1.GetFileByNameRequest
	(*GetFileByNameResponse)(nil),   // 4: files.v1.GetFileByNameResponse
	(*CreateFileRequest)(nil),       // 5: files.v1.CreateFileRequest
	(*CreateFileResponse)(nil),      // 6: files.v1.CreateFileResponse
	(*UpdateFileRequest)(nil),       // 7: files.v1.UpdateFileRequest
	(*UpdateFileResponse)(nil),      // 8: files.v1.UpdateFileResponse
	(*DeleteFileRequest)(nil),       // 9: files.v1.DeleteFileRequest
	(*DeleteFileResponse)(nil),      // 10: files.v1.DeleteFileResponse
}
var file_api_proto_files_v1_files_proto_depIdxs = []int32{
	0,  // 0: files.v1.GetFileByUserIdRequest.file:type_name -> files.v1.File
	0,  // 1: files.v1.GetFileByUserIdResponse.file:type_name -> files.v1.File
	0,  // 2: files.v1.GetFileByNameRequest.file:type_name -> files.v1.File
	0,  // 3: files.v1.GetFileByNameResponse.file:type_name -> files.v1.File
	0,  // 4: files.v1.CreateFileRequest.file:type_name -> files.v1.File
	0,  // 5: files.v1.UpdateFileRequest.file:type_name -> files.v1.File
	0,  // 6: files.v1.DeleteFileRequest.file:type_name -> files.v1.File
	5,  // 7: files.v1.FilesService.CreateFile:input_type -> files.v1.CreateFileRequest
	7,  // 8: files.v1.FilesService.UpdateFile:input_type -> files.v1.UpdateFileRequest
	9,  // 9: files.v1.FilesService.DeleteFile:input_type -> files.v1.DeleteFileRequest
	1,  // 10: files.v1.FilesService.GetFileByUserId:input_type -> files.v1.GetFileByUserIdRequest
	3,  // 11: files.v1.FilesService.GetFileByName:input_type -> files.v1.GetFileByNameRequest
	6,  // 12: files.v1.FilesService.CreateFile:output_type -> files.v1.CreateFileResponse
	8,  // 13: files.v1.FilesService.UpdateFile:output_type -> files.v1.UpdateFileResponse
	10, // 14: files.v1.FilesService.DeleteFile:output_type -> files.v1.DeleteFileResponse
	2,  // 15: files.v1.FilesService.GetFileByUserId:output_type -> files.v1.GetFileByUserIdResponse
	4,  // 16: files.v1.FilesService.GetFileByName:output_type -> files.v1.GetFileByNameResponse
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_api_proto_files_v1_files_proto_init() }
func file_api_proto_files_v1_files_proto_init() {
	if File_api_proto_files_v1_files_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_files_v1_files_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_files_v1_files_proto_goTypes,
		DependencyIndexes: file_api_proto_files_v1_files_proto_depIdxs,
		MessageInfos:      file_api_proto_files_v1_files_proto_msgTypes,
	}.Build()
	File_api_proto_files_v1_files_proto = out.File
	file_api_proto_files_v1_files_proto_rawDesc = nil
	file_api_proto_files_v1_files_proto_goTypes = nil
	file_api_proto_files_v1_files_proto_depIdxs = nil
}
