// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: darkchess/action_service.proto

package darkchessgrpc

import (
	_ "card-game-server-prototype/pkg/grpc"
	commongrpc "card-game-server-prototype/pkg/grpc/commongrpc"
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

type SkipScoreboardRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SkipScoreboardRequest) Reset() {
	*x = SkipScoreboardRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SkipScoreboardRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SkipScoreboardRequest) ProtoMessage() {}

func (x *SkipScoreboardRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SkipScoreboardRequest.ProtoReflect.Descriptor instead.
func (*SkipScoreboardRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{0}
}

type UpdatePlaySettingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsAuto bool `protobuf:"varint,1,opt,name=is_auto,json=isAuto,proto3" json:"is_auto,omitempty"`
}

func (x *UpdatePlaySettingRequest) Reset() {
	*x = UpdatePlaySettingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePlaySettingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePlaySettingRequest) ProtoMessage() {}

func (x *UpdatePlaySettingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePlaySettingRequest.ProtoReflect.Descriptor instead.
func (*UpdatePlaySettingRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{1}
}

func (x *UpdatePlaySettingRequest) GetIsAuto() bool {
	if x != nil {
		return x.IsAuto
	}
	return false
}

type PickRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PickRequest) Reset() {
	*x = PickRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PickRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PickRequest) ProtoMessage() {}

func (x *PickRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PickRequest.ProtoReflect.Descriptor instead.
func (*PickRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{2}
}

type RevealRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GridPosition *GridPosition `protobuf:"bytes,1,opt,name=grid_position,json=gridPosition,proto3" json:"grid_position,omitempty"`
}

func (x *RevealRequest) Reset() {
	*x = RevealRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RevealRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevealRequest) ProtoMessage() {}

func (x *RevealRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevealRequest.ProtoReflect.Descriptor instead.
func (*RevealRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{3}
}

func (x *RevealRequest) GetGridPosition() *GridPosition {
	if x != nil {
		return x.GridPosition
	}
	return nil
}

type MoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoTo      *GridPosition           `protobuf:"bytes,1,opt,name=go_to,json=goTo,proto3" json:"go_to,omitempty"`
	MovePiece commongrpc.CnChessPiece `protobuf:"varint,2,opt,name=move_piece,json=movePiece,proto3,enum=common.CnChessPiece" json:"move_piece,omitempty"`
}

func (x *MoveRequest) Reset() {
	*x = MoveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MoveRequest) ProtoMessage() {}

func (x *MoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MoveRequest.ProtoReflect.Descriptor instead.
func (*MoveRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{4}
}

func (x *MoveRequest) GetGoTo() *GridPosition {
	if x != nil {
		return x.GoTo
	}
	return nil
}

func (x *MoveRequest) GetMovePiece() commongrpc.CnChessPiece {
	if x != nil {
		return x.MovePiece
	}
	return commongrpc.CnChessPiece(0)
}

type CaptureRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoTo          *GridPosition           `protobuf:"bytes,1,opt,name=go_to,json=goTo,proto3" json:"go_to,omitempty"`
	MovePiece     commongrpc.CnChessPiece `protobuf:"varint,2,opt,name=move_piece,json=movePiece,proto3,enum=common.CnChessPiece" json:"move_piece,omitempty"`
	CapturedPiece commongrpc.CnChessPiece `protobuf:"varint,3,opt,name=captured_piece,json=capturedPiece,proto3,enum=common.CnChessPiece" json:"captured_piece,omitempty"`
}

func (x *CaptureRequest) Reset() {
	*x = CaptureRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CaptureRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaptureRequest) ProtoMessage() {}

func (x *CaptureRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaptureRequest.ProtoReflect.Descriptor instead.
func (*CaptureRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{5}
}

func (x *CaptureRequest) GetGoTo() *GridPosition {
	if x != nil {
		return x.GoTo
	}
	return nil
}

func (x *CaptureRequest) GetMovePiece() commongrpc.CnChessPiece {
	if x != nil {
		return x.MovePiece
	}
	return commongrpc.CnChessPiece(0)
}

func (x *CaptureRequest) GetCapturedPiece() commongrpc.CnChessPiece {
	if x != nil {
		return x.CapturedPiece
	}
	return commongrpc.CnChessPiece(0)
}

type SurrenderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *SurrenderRequest) Reset() {
	*x = SurrenderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SurrenderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SurrenderRequest) ProtoMessage() {}

func (x *SurrenderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SurrenderRequest.ProtoReflect.Descriptor instead.
func (*SurrenderRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{6}
}

func (x *SurrenderRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

type ClaimDrawRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"` // act uid
}

func (x *ClaimDrawRequest) Reset() {
	*x = ClaimDrawRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClaimDrawRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClaimDrawRequest) ProtoMessage() {}

func (x *ClaimDrawRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClaimDrawRequest.ProtoReflect.Descriptor instead.
func (*ClaimDrawRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{7}
}

func (x *ClaimDrawRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

type AnswerDrawRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid      string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"` // act uid
	IsAccept bool   `protobuf:"varint,2,opt,name=is_accept,json=isAccept,proto3" json:"is_accept,omitempty"`
}

func (x *AnswerDrawRequest) Reset() {
	*x = AnswerDrawRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnswerDrawRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnswerDrawRequest) ProtoMessage() {}

func (x *AnswerDrawRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnswerDrawRequest.ProtoReflect.Descriptor instead.
func (*AnswerDrawRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{8}
}

func (x *AnswerDrawRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *AnswerDrawRequest) GetIsAccept() bool {
	if x != nil {
		return x.IsAccept
	}
	return false
}

type AskExtraSecondsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
}

func (x *AskExtraSecondsRequest) Reset() {
	*x = AskExtraSecondsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_action_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AskExtraSecondsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AskExtraSecondsRequest) ProtoMessage() {}

func (x *AskExtraSecondsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_action_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AskExtraSecondsRequest.ProtoReflect.Descriptor instead.
func (*AskExtraSecondsRequest) Descriptor() ([]byte, []int) {
	return file_darkchess_action_service_proto_rawDescGZIP(), []int{9}
}

func (x *AskExtraSecondsRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

var File_darkchess_action_service_proto protoreflect.FileDescriptor

var file_darkchess_action_service_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2f, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x1a, 0x13, 0x67, 0x6c, 0x6f,
	0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x72, 0x65,
	0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x63, 0x68, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x69, 0x65, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2f,
	0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x17, 0x0a, 0x15, 0x53,
	0x6b, 0x69, 0x70, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x33, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6c,
	0x61, 0x79, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x61, 0x75, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x69, 0x73, 0x41, 0x75, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x50, 0x69, 0x63,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4d, 0x0a, 0x0d, 0x52, 0x65, 0x76, 0x65,
	0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3c, 0x0a, 0x0d, 0x67, 0x72, 0x69,
	0x64, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x47, 0x72, 0x69,
	0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x67, 0x72, 0x69, 0x64, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x70, 0x0a, 0x0b, 0x4d, 0x6f, 0x76, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x05, 0x67, 0x6f, 0x5f, 0x74, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73,
	0x73, 0x2e, 0x47, 0x72, 0x69, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04,
	0x67, 0x6f, 0x54, 0x6f, 0x12, 0x33, 0x0a, 0x0a, 0x6d, 0x6f, 0x76, 0x65, 0x5f, 0x70, 0x69, 0x65,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x43, 0x6e, 0x43, 0x68, 0x65, 0x73, 0x73, 0x50, 0x69, 0x65, 0x63, 0x65, 0x52, 0x09,
	0x6d, 0x6f, 0x76, 0x65, 0x50, 0x69, 0x65, 0x63, 0x65, 0x22, 0xb0, 0x01, 0x0a, 0x0e, 0x43, 0x61,
	0x70, 0x74, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x05,
	0x67, 0x6f, 0x5f, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x61,
	0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x47, 0x72, 0x69, 0x64, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x67, 0x6f, 0x54, 0x6f, 0x12, 0x33, 0x0a, 0x0a, 0x6d, 0x6f,
	0x76, 0x65, 0x5f, 0x70, 0x69, 0x65, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x6e, 0x43, 0x68, 0x65, 0x73, 0x73, 0x50,
	0x69, 0x65, 0x63, 0x65, 0x52, 0x09, 0x6d, 0x6f, 0x76, 0x65, 0x50, 0x69, 0x65, 0x63, 0x65, 0x12,
	0x3b, 0x0a, 0x0e, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x64, 0x5f, 0x70, 0x69, 0x65, 0x63,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x43, 0x6e, 0x43, 0x68, 0x65, 0x73, 0x73, 0x50, 0x69, 0x65, 0x63, 0x65, 0x52, 0x0d, 0x63,
	0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x64, 0x50, 0x69, 0x65, 0x63, 0x65, 0x22, 0x24, 0x0a, 0x10,
	0x53, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x69, 0x64, 0x22, 0x24, 0x0a, 0x10, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x44, 0x72, 0x61, 0x77, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x22, 0x42, 0x0a, 0x11, 0x41, 0x6e, 0x73, 0x77,
	0x65, 0x72, 0x44, 0x72, 0x61, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x22, 0x2a, 0x0a, 0x16,
	0x41, 0x73, 0x6b, 0x45, 0x78, 0x74, 0x72, 0x61, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x32, 0xe6, 0x07, 0x0a, 0x0d, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x06, 0x52, 0x65,
	0x73, 0x79, 0x6e, 0x63, 0x12, 0x15, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65,
	0x73, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x04, 0x4b, 0x69, 0x63, 0x6b, 0x12, 0x13, 0x2e,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4b, 0x69, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x05,
	0x41, 0x64, 0x64, 0x41, 0x69, 0x12, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x41,
	0x64, 0x64, 0x41, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x05, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x14,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3f,
	0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12,
	0x4c, 0x0a, 0x0e, 0x53, 0x6b, 0x69, 0x70, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x62, 0x6f, 0x61, 0x72,
	0x64, 0x12, 0x20, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x53, 0x6b,
	0x69, 0x70, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x52, 0x0a,
	0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x12, 0x23, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x38, 0x0a, 0x04, 0x50, 0x69, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x64, 0x61, 0x72, 0x6b,
	0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x50, 0x69, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x06, 0x52,
	0x65, 0x76, 0x65, 0x61, 0x6c, 0x12, 0x18, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73,
	0x73, 0x2e, 0x52, 0x65, 0x76, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x04, 0x4d, 0x6f, 0x76,
	0x65, 0x12, 0x16, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x4d, 0x6f,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x07, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x12, 0x19,
	0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x43, 0x61, 0x70, 0x74, 0x75,
	0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x09, 0x53, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x12, 0x1b, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x53, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x09, 0x43, 0x6c, 0x61, 0x69, 0x6d,
	0x44, 0x72, 0x61, 0x77, 0x12, 0x1b, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73,
	0x2e, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x44, 0x72, 0x61, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0a, 0x41,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x44, 0x72, 0x61, 0x77, 0x12, 0x1c, 0x2e, 0x64, 0x61, 0x72, 0x6b,
	0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x44, 0x72, 0x61, 0x77,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x4e, 0x0a, 0x0f, 0x41, 0x73, 0x6b, 0x45, 0x78, 0x74, 0x72, 0x61, 0x53, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x73, 0x12, 0x21, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73,
	0x2e, 0x41, 0x73, 0x6b, 0x45, 0x78, 0x74, 0x72, 0x61, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x42, 0x74, 0x5a, 0x32, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65,
	0x73, 0x73, 0x67, 0x72, 0x70, 0x63, 0x3b, 0xaa, 0x02, 0x21, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x44, 0x61, 0x72, 0x6b, 0x43, 0x68, 0x65,
	0x73, 0x73, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0xca, 0xb2, 0x04, 0x18, 0x4a,
	0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x44, 0x61,
	0x72, 0x6b, 0x43, 0x68, 0x65, 0x73, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_darkchess_action_service_proto_rawDescOnce sync.Once
	file_darkchess_action_service_proto_rawDescData = file_darkchess_action_service_proto_rawDesc
)

func file_darkchess_action_service_proto_rawDescGZIP() []byte {
	file_darkchess_action_service_proto_rawDescOnce.Do(func() {
		file_darkchess_action_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_darkchess_action_service_proto_rawDescData)
	})
	return file_darkchess_action_service_proto_rawDescData
}

var file_darkchess_action_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_darkchess_action_service_proto_goTypes = []interface{}{
	(*SkipScoreboardRequest)(nil),       // 0: darkchess.SkipScoreboardRequest
	(*UpdatePlaySettingRequest)(nil),    // 1: darkchess.UpdatePlaySettingRequest
	(*PickRequest)(nil),                 // 2: darkchess.PickRequest
	(*RevealRequest)(nil),               // 3: darkchess.RevealRequest
	(*MoveRequest)(nil),                 // 4: darkchess.MoveRequest
	(*CaptureRequest)(nil),              // 5: darkchess.CaptureRequest
	(*SurrenderRequest)(nil),            // 6: darkchess.SurrenderRequest
	(*ClaimDrawRequest)(nil),            // 7: darkchess.ClaimDrawRequest
	(*AnswerDrawRequest)(nil),           // 8: darkchess.AnswerDrawRequest
	(*AskExtraSecondsRequest)(nil),      // 9: darkchess.AskExtraSecondsRequest
	(*GridPosition)(nil),                // 10: darkchess.GridPosition
	(commongrpc.CnChessPiece)(0),        // 11: common.CnChessPiece
	(*commongrpc.ResyncRequest)(nil),    // 12: common.ResyncRequest
	(*commongrpc.KickRequest)(nil),      // 13: common.KickRequest
	(*commongrpc.AddAiRequest)(nil),     // 14: common.AddAiRequest
	(*commongrpc.ReadyRequest)(nil),     // 15: common.ReadyRequest
	(*commongrpc.StartGameRequest)(nil), // 16: common.StartGameRequest
	(*emptypb.Empty)(nil),               // 17: google.protobuf.Empty
}
var file_darkchess_action_service_proto_depIdxs = []int32{
	10, // 0: darkchess.RevealRequest.grid_position:type_name -> darkchess.GridPosition
	10, // 1: darkchess.MoveRequest.go_to:type_name -> darkchess.GridPosition
	11, // 2: darkchess.MoveRequest.move_piece:type_name -> common.CnChessPiece
	10, // 3: darkchess.CaptureRequest.go_to:type_name -> darkchess.GridPosition
	11, // 4: darkchess.CaptureRequest.move_piece:type_name -> common.CnChessPiece
	11, // 5: darkchess.CaptureRequest.captured_piece:type_name -> common.CnChessPiece
	12, // 6: darkchess.ActionService.Resync:input_type -> common.ResyncRequest
	13, // 7: darkchess.ActionService.Kick:input_type -> common.KickRequest
	14, // 8: darkchess.ActionService.AddAi:input_type -> common.AddAiRequest
	15, // 9: darkchess.ActionService.Ready:input_type -> common.ReadyRequest
	16, // 10: darkchess.ActionService.StartGame:input_type -> common.StartGameRequest
	0,  // 11: darkchess.ActionService.SkipScoreboard:input_type -> darkchess.SkipScoreboardRequest
	1,  // 12: darkchess.ActionService.UpdatePlaySetting:input_type -> darkchess.UpdatePlaySettingRequest
	2,  // 13: darkchess.ActionService.Pick:input_type -> darkchess.PickRequest
	3,  // 14: darkchess.ActionService.Reveal:input_type -> darkchess.RevealRequest
	4,  // 15: darkchess.ActionService.Move:input_type -> darkchess.MoveRequest
	5,  // 16: darkchess.ActionService.Capture:input_type -> darkchess.CaptureRequest
	6,  // 17: darkchess.ActionService.Surrender:input_type -> darkchess.SurrenderRequest
	7,  // 18: darkchess.ActionService.ClaimDraw:input_type -> darkchess.ClaimDrawRequest
	8,  // 19: darkchess.ActionService.AnswerDraw:input_type -> darkchess.AnswerDrawRequest
	9,  // 20: darkchess.ActionService.AskExtraSeconds:input_type -> darkchess.AskExtraSecondsRequest
	17, // 21: darkchess.ActionService.Resync:output_type -> google.protobuf.Empty
	17, // 22: darkchess.ActionService.Kick:output_type -> google.protobuf.Empty
	17, // 23: darkchess.ActionService.AddAi:output_type -> google.protobuf.Empty
	17, // 24: darkchess.ActionService.Ready:output_type -> google.protobuf.Empty
	17, // 25: darkchess.ActionService.StartGame:output_type -> google.protobuf.Empty
	17, // 26: darkchess.ActionService.SkipScoreboard:output_type -> google.protobuf.Empty
	17, // 27: darkchess.ActionService.UpdatePlaySetting:output_type -> google.protobuf.Empty
	17, // 28: darkchess.ActionService.Pick:output_type -> google.protobuf.Empty
	17, // 29: darkchess.ActionService.Reveal:output_type -> google.protobuf.Empty
	17, // 30: darkchess.ActionService.Move:output_type -> google.protobuf.Empty
	17, // 31: darkchess.ActionService.Capture:output_type -> google.protobuf.Empty
	17, // 32: darkchess.ActionService.Surrender:output_type -> google.protobuf.Empty
	17, // 33: darkchess.ActionService.ClaimDraw:output_type -> google.protobuf.Empty
	17, // 34: darkchess.ActionService.AnswerDraw:output_type -> google.protobuf.Empty
	17, // 35: darkchess.ActionService.AskExtraSeconds:output_type -> google.protobuf.Empty
	21, // [21:36] is the sub-list for method output_type
	6,  // [6:21] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_darkchess_action_service_proto_init() }
func file_darkchess_action_service_proto_init() {
	if File_darkchess_action_service_proto != nil {
		return
	}
	file_darkchess_board_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_darkchess_action_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SkipScoreboardRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePlaySettingRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PickRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RevealRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MoveRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CaptureRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SurrenderRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClaimDrawRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnswerDrawRequest); i {
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
		file_darkchess_action_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AskExtraSecondsRequest); i {
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
			RawDescriptor: file_darkchess_action_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_darkchess_action_service_proto_goTypes,
		DependencyIndexes: file_darkchess_action_service_proto_depIdxs,
		MessageInfos:      file_darkchess_action_service_proto_msgTypes,
	}.Build()
	File_darkchess_action_service_proto = out.File
	file_darkchess_action_service_proto_rawDesc = nil
	file_darkchess_action_service_proto_goTypes = nil
	file_darkchess_action_service_proto_depIdxs = nil
}
