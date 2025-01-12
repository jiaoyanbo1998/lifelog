// 版本

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: api/proto/comment/v1/comment.proto

// 包名

package commentv1

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

// 评论信息结构体
type Comment struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Biz           string                 `protobuf:"bytes,2,opt,name=biz,proto3" json:"biz,omitempty"`
	BizId         int64                  `protobuf:"varint,3,opt,name=biz_id,json=bizId,proto3" json:"biz_id,omitempty"`
	Content       string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	ParentId      int64                  `protobuf:"varint,5,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`
	RootId        int64                  `protobuf:"varint,6,opt,name=root_id,json=rootId,proto3" json:"root_id,omitempty"`
	Id            int64                  `protobuf:"varint,7,opt,name=id,proto3" json:"id,omitempty"`
	CreateTime    int64                  `protobuf:"varint,8,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	UpdateTime    int64                  `protobuf:"varint,9,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	Uuid          string                 `protobuf:"bytes,10,opt,name=uuid,proto3" json:"uuid,omitempty"`
	TargetUserId  int64                  `protobuf:"varint,11,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Comment) Reset() {
	*x = Comment{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Comment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Comment) ProtoMessage() {}

func (x *Comment) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Comment.ProtoReflect.Descriptor instead.
func (*Comment) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{0}
}

func (x *Comment) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Comment) GetBiz() string {
	if x != nil {
		return x.Biz
	}
	return ""
}

func (x *Comment) GetBizId() int64 {
	if x != nil {
		return x.BizId
	}
	return 0
}

func (x *Comment) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Comment) GetParentId() int64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *Comment) GetRootId() int64 {
	if x != nil {
		return x.RootId
	}
	return 0
}

func (x *Comment) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Comment) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *Comment) GetUpdateTime() int64 {
	if x != nil {
		return x.UpdateTime
	}
	return 0
}

func (x *Comment) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *Comment) GetTargetUserId() int64 {
	if x != nil {
		return x.TargetUserId
	}
	return 0
}

// 批量创建评论
// 请求
type ProducerCommentEventRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comment       *Comment               `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProducerCommentEventRequest) Reset() {
	*x = ProducerCommentEventRequest{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProducerCommentEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProducerCommentEventRequest) ProtoMessage() {}

func (x *ProducerCommentEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProducerCommentEventRequest.ProtoReflect.Descriptor instead.
func (*ProducerCommentEventRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{1}
}

func (x *ProducerCommentEventRequest) GetComment() *Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

// 响应
type ProducerCommentEventResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProducerCommentEventResponse) Reset() {
	*x = ProducerCommentEventResponse{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProducerCommentEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProducerCommentEventResponse) ProtoMessage() {}

func (x *ProducerCommentEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProducerCommentEventResponse.ProtoReflect.Descriptor instead.
func (*ProducerCommentEventResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{2}
}

// 删除评论
// 请求
type DeleteCommentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCommentRequest) Reset() {
	*x = DeleteCommentRequest{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCommentRequest) ProtoMessage() {}

func (x *DeleteCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCommentRequest.ProtoReflect.Descriptor instead.
func (*DeleteCommentRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteCommentRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// 响应
type DeleteCommentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCommentResponse) Reset() {
	*x = DeleteCommentResponse{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCommentResponse) ProtoMessage() {}

func (x *DeleteCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCommentResponse.ProtoReflect.Descriptor instead.
func (*DeleteCommentResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{4}
}

// 编辑评论
// 请求
type EditCommentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comment       *Comment               `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EditCommentRequest) Reset() {
	*x = EditCommentRequest{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EditCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditCommentRequest) ProtoMessage() {}

func (x *EditCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditCommentRequest.ProtoReflect.Descriptor instead.
func (*EditCommentRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{5}
}

func (x *EditCommentRequest) GetComment() *Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

// 响应
type EditCommentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EditCommentResponse) Reset() {
	*x = EditCommentResponse{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EditCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditCommentResponse) ProtoMessage() {}

func (x *EditCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditCommentResponse.ProtoReflect.Descriptor instead.
func (*EditCommentResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{6}
}

// 查找一级评论（parent_id==null）
// 请求
type FirstListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comment       *Comment               `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	Min           int64                  `protobuf:"varint,2,opt,name=min,proto3" json:"min,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FirstListRequest) Reset() {
	*x = FirstListRequest{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FirstListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FirstListRequest) ProtoMessage() {}

func (x *FirstListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FirstListRequest.ProtoReflect.Descriptor instead.
func (*FirstListRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{7}
}

func (x *FirstListRequest) GetComment() *Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

func (x *FirstListRequest) GetMin() int64 {
	if x != nil {
		return x.Min
	}
	return 0
}

// 响应
type FirstListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comments      []*Comment             `protobuf:"bytes,1,rep,name=comments,proto3" json:"comments,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FirstListResponse) Reset() {
	*x = FirstListResponse{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FirstListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FirstListResponse) ProtoMessage() {}

func (x *FirstListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FirstListResponse.ProtoReflect.Descriptor instead.
func (*FirstListResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{8}
}

func (x *FirstListResponse) GetComments() []*Comment {
	if x != nil {
		return x.Comments
	}
	return nil
}

// 查找根评论下的所有子孙评论，root_id = 根评论id，根据id降序排序（先发表的评论id比后发表的小）
// 请求
type EveryRootChildSonListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comment       *Comment               `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	Limit         int64                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"` // 限制条数
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EveryRootChildSonListRequest) Reset() {
	*x = EveryRootChildSonListRequest{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EveryRootChildSonListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EveryRootChildSonListRequest) ProtoMessage() {}

func (x *EveryRootChildSonListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EveryRootChildSonListRequest.ProtoReflect.Descriptor instead.
func (*EveryRootChildSonListRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{9}
}

func (x *EveryRootChildSonListRequest) GetComment() *Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

func (x *EveryRootChildSonListRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// 响应
type EveryRootChildSonListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comments      []*Comment             `protobuf:"bytes,1,rep,name=comments,proto3" json:"comments,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EveryRootChildSonListResponse) Reset() {
	*x = EveryRootChildSonListResponse{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EveryRootChildSonListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EveryRootChildSonListResponse) ProtoMessage() {}

func (x *EveryRootChildSonListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EveryRootChildSonListResponse.ProtoReflect.Descriptor instead.
func (*EveryRootChildSonListResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{10}
}

func (x *EveryRootChildSonListResponse) GetComments() []*Comment {
	if x != nil {
		return x.Comments
	}
	return nil
}

// 查询某个评论的，一级子评论
// 请求
type SonListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comment       *Comment               `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	Limit         int64                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"` // 限制条数
	Offset        int64                  `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SonListRequest) Reset() {
	*x = SonListRequest{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SonListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SonListRequest) ProtoMessage() {}

func (x *SonListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SonListRequest.ProtoReflect.Descriptor instead.
func (*SonListRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{11}
}

func (x *SonListRequest) GetComment() *Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

func (x *SonListRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *SonListRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

// 响应
type SonListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comments      []*Comment             `protobuf:"bytes,1,rep,name=comments,proto3" json:"comments,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SonListResponse) Reset() {
	*x = SonListResponse{}
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SonListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SonListResponse) ProtoMessage() {}

func (x *SonListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_comment_v1_comment_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SonListResponse.ProtoReflect.Descriptor instead.
func (*SonListResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_comment_v1_comment_proto_rawDescGZIP(), []int{12}
}

func (x *SonListResponse) GetComments() []*Comment {
	if x != nil {
		return x.Comments
	}
	return nil
}

var File_api_proto_comment_v1_comment_proto protoreflect.FileDescriptor

var file_api_proto_comment_v1_comment_proto_rawDesc = []byte{
	0x0a, 0x22, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x22, 0xa7, 0x02, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x69, 0x7a, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x62, 0x69, 0x7a, 0x12, 0x15, 0x0a, 0x06, 0x62, 0x69, 0x7a, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x62, 0x69, 0x7a, 0x49, 0x64, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x61, 0x72,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6f, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x74, 0x49, 0x64, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x75, 0x75, 0x69, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4c, 0x0a, 0x1b, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x1e, 0x0a, 0x1c, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x26, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x17, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x43, 0x0a, 0x12, 0x45, 0x64, 0x69,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2d, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x15,
	0x0a, 0x13, 0x45, 0x64, 0x69, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x53, 0x0a, 0x10, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x22, 0x44, 0x0a, 0x11, 0x46, 0x69,
	0x72, 0x73, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2f, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x22, 0x63, 0x0a, 0x1c, 0x45, 0x76, 0x65, 0x72, 0x79, 0x52, 0x6f, 0x6f, 0x74, 0x43, 0x68, 0x69,
	0x6c, 0x64, 0x53, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2d, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x50, 0x0a, 0x1d, 0x45, 0x76, 0x65, 0x72, 0x79, 0x52, 0x6f,
	0x6f, 0x74, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x53, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x6d, 0x0a, 0x0e, 0x53, 0x6f, 0x6e, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x07, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x42, 0x0a, 0x0f, 0x53, 0x6f, 0x6e, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0x9d, 0x04, 0x0a, 0x0e, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a,
	0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x20,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x21, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x0b, 0x45, 0x64, 0x69, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x1e, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x45, 0x64, 0x69, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x45, 0x64, 0x69, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x09, 0x46, 0x69, 0x72, 0x73, 0x74, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x1c, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69,
	0x72, 0x73, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x72, 0x73,
	0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6c, 0x0a,
	0x15, 0x45, 0x76, 0x65, 0x72, 0x79, 0x52, 0x6f, 0x6f, 0x74, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x53,
	0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x72, 0x79, 0x52, 0x6f, 0x6f, 0x74, 0x43, 0x68, 0x69,
	0x6c, 0x64, 0x53, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x29, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76,
	0x65, 0x72, 0x79, 0x52, 0x6f, 0x6f, 0x74, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x53, 0x6f, 0x6e, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x07, 0x53,
	0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1a, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x69, 0x0a, 0x14, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x27, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x28, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x95, 0x01, 0x0a, 0x0e, 0x63,
	0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2c, 0x61,
	0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76,
	0x31, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x58,
	0x58, 0xaa, 0x02, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x16, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_comment_v1_comment_proto_rawDescOnce sync.Once
	file_api_proto_comment_v1_comment_proto_rawDescData = file_api_proto_comment_v1_comment_proto_rawDesc
)

func file_api_proto_comment_v1_comment_proto_rawDescGZIP() []byte {
	file_api_proto_comment_v1_comment_proto_rawDescOnce.Do(func() {
		file_api_proto_comment_v1_comment_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_comment_v1_comment_proto_rawDescData)
	})
	return file_api_proto_comment_v1_comment_proto_rawDescData
}

var file_api_proto_comment_v1_comment_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_api_proto_comment_v1_comment_proto_goTypes = []any{
	(*Comment)(nil),                       // 0: comment.v1.Comment
	(*ProducerCommentEventRequest)(nil),   // 1: comment.v1.ProducerCommentEventRequest
	(*ProducerCommentEventResponse)(nil),  // 2: comment.v1.ProducerCommentEventResponse
	(*DeleteCommentRequest)(nil),          // 3: comment.v1.DeleteCommentRequest
	(*DeleteCommentResponse)(nil),         // 4: comment.v1.DeleteCommentResponse
	(*EditCommentRequest)(nil),            // 5: comment.v1.EditCommentRequest
	(*EditCommentResponse)(nil),           // 6: comment.v1.EditCommentResponse
	(*FirstListRequest)(nil),              // 7: comment.v1.FirstListRequest
	(*FirstListResponse)(nil),             // 8: comment.v1.FirstListResponse
	(*EveryRootChildSonListRequest)(nil),  // 9: comment.v1.EveryRootChildSonListRequest
	(*EveryRootChildSonListResponse)(nil), // 10: comment.v1.EveryRootChildSonListResponse
	(*SonListRequest)(nil),                // 11: comment.v1.SonListRequest
	(*SonListResponse)(nil),               // 12: comment.v1.SonListResponse
}
var file_api_proto_comment_v1_comment_proto_depIdxs = []int32{
	0,  // 0: comment.v1.ProducerCommentEventRequest.comment:type_name -> comment.v1.Comment
	0,  // 1: comment.v1.EditCommentRequest.comment:type_name -> comment.v1.Comment
	0,  // 2: comment.v1.FirstListRequest.comment:type_name -> comment.v1.Comment
	0,  // 3: comment.v1.FirstListResponse.comments:type_name -> comment.v1.Comment
	0,  // 4: comment.v1.EveryRootChildSonListRequest.comment:type_name -> comment.v1.Comment
	0,  // 5: comment.v1.EveryRootChildSonListResponse.comments:type_name -> comment.v1.Comment
	0,  // 6: comment.v1.SonListRequest.comment:type_name -> comment.v1.Comment
	0,  // 7: comment.v1.SonListResponse.comments:type_name -> comment.v1.Comment
	3,  // 8: comment.v1.CommentService.DeleteComment:input_type -> comment.v1.DeleteCommentRequest
	5,  // 9: comment.v1.CommentService.EditComment:input_type -> comment.v1.EditCommentRequest
	7,  // 10: comment.v1.CommentService.FirstList:input_type -> comment.v1.FirstListRequest
	9,  // 11: comment.v1.CommentService.EveryRootChildSonList:input_type -> comment.v1.EveryRootChildSonListRequest
	11, // 12: comment.v1.CommentService.SonList:input_type -> comment.v1.SonListRequest
	1,  // 13: comment.v1.CommentService.ProducerCommentEvent:input_type -> comment.v1.ProducerCommentEventRequest
	4,  // 14: comment.v1.CommentService.DeleteComment:output_type -> comment.v1.DeleteCommentResponse
	6,  // 15: comment.v1.CommentService.EditComment:output_type -> comment.v1.EditCommentResponse
	8,  // 16: comment.v1.CommentService.FirstList:output_type -> comment.v1.FirstListResponse
	10, // 17: comment.v1.CommentService.EveryRootChildSonList:output_type -> comment.v1.EveryRootChildSonListResponse
	12, // 18: comment.v1.CommentService.SonList:output_type -> comment.v1.SonListResponse
	2,  // 19: comment.v1.CommentService.ProducerCommentEvent:output_type -> comment.v1.ProducerCommentEventResponse
	14, // [14:20] is the sub-list for method output_type
	8,  // [8:14] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_api_proto_comment_v1_comment_proto_init() }
func file_api_proto_comment_v1_comment_proto_init() {
	if File_api_proto_comment_v1_comment_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_comment_v1_comment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_comment_v1_comment_proto_goTypes,
		DependencyIndexes: file_api_proto_comment_v1_comment_proto_depIdxs,
		MessageInfos:      file_api_proto_comment_v1_comment_proto_msgTypes,
	}.Build()
	File_api_proto_comment_v1_comment_proto = out.File
	file_api_proto_comment_v1_comment_proto_rawDesc = nil
	file_api_proto_comment_v1_comment_proto_goTypes = nil
	file_api_proto_comment_v1_comment_proto_depIdxs = nil
}
