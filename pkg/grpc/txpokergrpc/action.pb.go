// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: txpoker/action.proto

package txpokergrpc

import (
	_ "card-game-server-prototype/pkg/grpc"
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

type StandUpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StandUpRequest) Reset() {
	*x = StandUpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StandUpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StandUpRequest) ProtoMessage() {}

func (x *StandUpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StandUpRequest.ProtoReflect.Descriptor instead.
func (*StandUpRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{0}
}

type SitDownRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SeatId int32 `protobuf:"varint,1,opt,name=seat_id,json=seatId,proto3" json:"seat_id,omitempty"`
}

func (x *SitDownRequest) Reset() {
	*x = SitDownRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SitDownRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SitDownRequest) ProtoMessage() {}

func (x *SitDownRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SitDownRequest.ProtoReflect.Descriptor instead.
func (*SitDownRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{1}
}

func (x *SitDownRequest) GetSeatId() int32 {
	if x != nil {
		return x.SeatId
	}
	return 0
}

type BuyInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BuyInChip int32 `protobuf:"varint,1,opt,name=buy_in_chip,json=buyInChip,proto3" json:"buy_in_chip,omitempty"`
}

func (x *BuyInRequest) Reset() {
	*x = BuyInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuyInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuyInRequest) ProtoMessage() {}

func (x *BuyInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuyInRequest.ProtoReflect.Descriptor instead.
func (*BuyInRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{2}
}

func (x *BuyInRequest) GetBuyInChip() int32 {
	if x != nil {
		return x.BuyInChip
	}
	return 0
}

type SitOutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SitOutRequest) Reset() {
	*x = SitOutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SitOutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SitOutRequest) ProtoMessage() {}

func (x *SitOutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SitOutRequest.ProtoReflect.Descriptor instead.
func (*SitOutRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{3}
}

type TopUpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TopUpChip int32 `protobuf:"varint,1,opt,name=top_up_chip,json=topUpChip,proto3" json:"top_up_chip,omitempty"`
}

func (x *TopUpRequest) Reset() {
	*x = TopUpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopUpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopUpRequest) ProtoMessage() {}

func (x *TopUpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopUpRequest.ProtoReflect.Descriptor instead.
func (*TopUpRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{4}
}

func (x *TopUpRequest) GetTopUpChip() int32 {
	if x != nil {
		return x.TopUpChip
	}
	return 0
}

type FoldRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FoldRequest) Reset() {
	*x = FoldRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FoldRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FoldRequest) ProtoMessage() {}

func (x *FoldRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FoldRequest.ProtoReflect.Descriptor instead.
func (*FoldRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{5}
}

type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{6}
}

type BetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chip int32 `protobuf:"varint,1,opt,name=chip,proto3" json:"chip,omitempty"`
}

func (x *BetRequest) Reset() {
	*x = BetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BetRequest) ProtoMessage() {}

func (x *BetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BetRequest.ProtoReflect.Descriptor instead.
func (*BetRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{7}
}

func (x *BetRequest) GetChip() int32 {
	if x != nil {
		return x.Chip
	}
	return 0
}

type CallRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CallRequest) Reset() {
	*x = CallRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CallRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallRequest) ProtoMessage() {}

func (x *CallRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallRequest.ProtoReflect.Descriptor instead.
func (*CallRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{8}
}

type RaiseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chip int32 `protobuf:"varint,1,opt,name=chip,proto3" json:"chip,omitempty"`
}

func (x *RaiseRequest) Reset() {
	*x = RaiseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaiseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaiseRequest) ProtoMessage() {}

func (x *RaiseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaiseRequest.ProtoReflect.Descriptor instead.
func (*RaiseRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{9}
}

func (x *RaiseRequest) GetChip() int32 {
	if x != nil {
		return x.Chip
	}
	return 0
}

type AllInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AllInRequest) Reset() {
	*x = AllInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllInRequest) ProtoMessage() {}

func (x *AllInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllInRequest.ProtoReflect.Descriptor instead.
func (*AllInRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{10}
}

type ShowFoldRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShowFoldType int32 `protobuf:"varint,1,opt,name=show_fold_type,json=showFoldType,proto3" json:"show_fold_type,omitempty"` // 0b00, 0b01, 0b10, 0b11
}

func (x *ShowFoldRequest) Reset() {
	*x = ShowFoldRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShowFoldRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShowFoldRequest) ProtoMessage() {}

func (x *ShowFoldRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShowFoldRequest.ProtoReflect.Descriptor instead.
func (*ShowFoldRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{11}
}

func (x *ShowFoldRequest) GetShowFoldType() int32 {
	if x != nil {
		return x.ShowFoldType
	}
	return 0
}

type UpdateWaitBBSettingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WaitBB bool `protobuf:"varint,1,opt,name=waitBB,proto3" json:"waitBB,omitempty"`
}

func (x *UpdateWaitBBSettingRequest) Reset() {
	*x = UpdateWaitBBSettingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateWaitBBSettingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateWaitBBSettingRequest) ProtoMessage() {}

func (x *UpdateWaitBBSettingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateWaitBBSettingRequest.ProtoReflect.Descriptor instead.
func (*UpdateWaitBBSettingRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{12}
}

func (x *UpdateWaitBBSettingRequest) GetWaitBB() bool {
	if x != nil {
		return x.WaitBB
	}
	return false
}

type UpdateAutoTopUpSettingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AutoTopUp                 bool    `protobuf:"varint,1,opt,name=auto_top_up,json=autoTopUp,proto3" json:"auto_top_up,omitempty"`
	AutoTopUpThresholdPercent float64 `protobuf:"fixed64,2,opt,name=auto_top_up_threshold_percent,json=autoTopUpThresholdPercent,proto3" json:"auto_top_up_threshold_percent,omitempty"`
	AutoTopUpChipPercent      float64 `protobuf:"fixed64,3,opt,name=auto_top_up_chip_percent,json=autoTopUpChipPercent,proto3" json:"auto_top_up_chip_percent,omitempty"`
}

func (x *UpdateAutoTopUpSettingRequest) Reset() {
	*x = UpdateAutoTopUpSettingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAutoTopUpSettingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAutoTopUpSettingRequest) ProtoMessage() {}

func (x *UpdateAutoTopUpSettingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAutoTopUpSettingRequest.ProtoReflect.Descriptor instead.
func (*UpdateAutoTopUpSettingRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_proto_rawDescGZIP(), []int{13}
}

func (x *UpdateAutoTopUpSettingRequest) GetAutoTopUp() bool {
	if x != nil {
		return x.AutoTopUp
	}
	return false
}

func (x *UpdateAutoTopUpSettingRequest) GetAutoTopUpThresholdPercent() float64 {
	if x != nil {
		return x.AutoTopUpThresholdPercent
	}
	return 0
}

func (x *UpdateAutoTopUpSettingRequest) GetAutoTopUpChipPercent() float64 {
	if x != nil {
		return x.AutoTopUpChipPercent
	}
	return 0
}

var File_txpoker_action_proto protoreflect.FileDescriptor

var file_txpoker_action_proto_rawDesc = []byte{
	0x0a, 0x14, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x1a,
	0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x10, 0x0a, 0x0e, 0x53, 0x74, 0x61, 0x6e, 0x64, 0x55, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x29, 0x0a, 0x0e, 0x53, 0x69, 0x74, 0x44, 0x6f, 0x77,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x65, 0x61, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x65, 0x61, 0x74, 0x49,
	0x64, 0x22, 0x2e, 0x0a, 0x0c, 0x42, 0x75, 0x79, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x0b, 0x62, 0x75, 0x79, 0x5f, 0x69, 0x6e, 0x5f, 0x63, 0x68, 0x69, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x62, 0x75, 0x79, 0x49, 0x6e, 0x43, 0x68, 0x69,
	0x70, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x69, 0x74, 0x4f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x2e, 0x0a, 0x0c, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0b, 0x74, 0x6f, 0x70, 0x5f, 0x75, 0x70, 0x5f, 0x63, 0x68, 0x69,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x6f, 0x70, 0x55, 0x70, 0x43, 0x68,
	0x69, 0x70, 0x22, 0x0d, 0x0a, 0x0b, 0x46, 0x6f, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x0e, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x20, 0x0a, 0x0a, 0x42, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x63, 0x68, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63,
	0x68, 0x69, 0x70, 0x22, 0x0d, 0x0a, 0x0b, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x22, 0x0a, 0x0c, 0x52, 0x61, 0x69, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x68, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x63, 0x68, 0x69, 0x70, 0x22, 0x0e, 0x0a, 0x0c, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x37, 0x0a, 0x0f, 0x53, 0x68, 0x6f, 0x77, 0x46, 0x6f,
	0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x68, 0x6f,
	0x77, 0x5f, 0x66, 0x6f, 0x6c, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0c, 0x73, 0x68, 0x6f, 0x77, 0x46, 0x6f, 0x6c, 0x64, 0x54, 0x79, 0x70, 0x65, 0x22,
	0x34, 0x0a, 0x1a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x61, 0x69, 0x74, 0x42, 0x42, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x77, 0x61, 0x69, 0x74, 0x42, 0x42, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x77,
	0x61, 0x69, 0x74, 0x42, 0x42, 0x22, 0xb9, 0x01, 0x0a, 0x1d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x41, 0x75, 0x74, 0x6f, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x6f, 0x5f,
	0x74, 0x6f, 0x70, 0x5f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x61, 0x75,
	0x74, 0x6f, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x12, 0x40, 0x0a, 0x1d, 0x61, 0x75, 0x74, 0x6f, 0x5f,
	0x74, 0x6f, 0x70, 0x5f, 0x75, 0x70, 0x5f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64,
	0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x19,
	0x61, 0x75, 0x74, 0x6f, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f,
	0x6c, 0x64, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x12, 0x36, 0x0a, 0x18, 0x61, 0x75, 0x74,
	0x6f, 0x5f, 0x74, 0x6f, 0x70, 0x5f, 0x75, 0x70, 0x5f, 0x63, 0x68, 0x69, 0x70, 0x5f, 0x70, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x14, 0x61, 0x75, 0x74,
	0x6f, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x43, 0x68, 0x69, 0x70, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e,
	0x74, 0x42, 0x6b, 0x5a, 0x2f, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72,
	0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1d, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d,
	0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x4d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0xca, 0xb2, 0x04, 0x16, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61,
	0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_txpoker_action_proto_rawDescOnce sync.Once
	file_txpoker_action_proto_rawDescData = file_txpoker_action_proto_rawDesc
)

func file_txpoker_action_proto_rawDescGZIP() []byte {
	file_txpoker_action_proto_rawDescOnce.Do(func() {
		file_txpoker_action_proto_rawDescData = protoimpl.X.CompressGZIP(file_txpoker_action_proto_rawDescData)
	})
	return file_txpoker_action_proto_rawDescData
}

var file_txpoker_action_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_txpoker_action_proto_goTypes = []interface{}{
	(*StandUpRequest)(nil),                // 0: txpoker.StandUpRequest
	(*SitDownRequest)(nil),                // 1: txpoker.SitDownRequest
	(*BuyInRequest)(nil),                  // 2: txpoker.BuyInRequest
	(*SitOutRequest)(nil),                 // 3: txpoker.SitOutRequest
	(*TopUpRequest)(nil),                  // 4: txpoker.TopUpRequest
	(*FoldRequest)(nil),                   // 5: txpoker.FoldRequest
	(*CheckRequest)(nil),                  // 6: txpoker.CheckRequest
	(*BetRequest)(nil),                    // 7: txpoker.BetRequest
	(*CallRequest)(nil),                   // 8: txpoker.CallRequest
	(*RaiseRequest)(nil),                  // 9: txpoker.RaiseRequest
	(*AllInRequest)(nil),                  // 10: txpoker.AllInRequest
	(*ShowFoldRequest)(nil),               // 11: txpoker.ShowFoldRequest
	(*UpdateWaitBBSettingRequest)(nil),    // 12: txpoker.UpdateWaitBBSettingRequest
	(*UpdateAutoTopUpSettingRequest)(nil), // 13: txpoker.UpdateAutoTopUpSettingRequest
}
var file_txpoker_action_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_txpoker_action_proto_init() }
func file_txpoker_action_proto_init() {
	if File_txpoker_action_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_txpoker_action_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StandUpRequest); i {
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
		file_txpoker_action_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SitDownRequest); i {
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
		file_txpoker_action_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuyInRequest); i {
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
		file_txpoker_action_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SitOutRequest); i {
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
		file_txpoker_action_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopUpRequest); i {
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
		file_txpoker_action_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FoldRequest); i {
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
		file_txpoker_action_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
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
		file_txpoker_action_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BetRequest); i {
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
		file_txpoker_action_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CallRequest); i {
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
		file_txpoker_action_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RaiseRequest); i {
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
		file_txpoker_action_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllInRequest); i {
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
		file_txpoker_action_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShowFoldRequest); i {
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
		file_txpoker_action_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateWaitBBSettingRequest); i {
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
		file_txpoker_action_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateAutoTopUpSettingRequest); i {
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
			RawDescriptor: file_txpoker_action_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_txpoker_action_proto_goTypes,
		DependencyIndexes: file_txpoker_action_proto_depIdxs,
		MessageInfos:      file_txpoker_action_proto_msgTypes,
	}.Build()
	File_txpoker_action_proto = out.File
	file_txpoker_action_proto_rawDesc = nil
	file_txpoker_action_proto_goTypes = nil
	file_txpoker_action_proto_depIdxs = nil
}
