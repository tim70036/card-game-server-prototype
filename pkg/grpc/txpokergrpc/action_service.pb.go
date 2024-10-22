// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: txpoker/action_service.proto

package txpokergrpc

import (
	_ "card-game-server-prototype/pkg/grpc"
	commongrpc "card-game-server-prototype/pkg/grpc/commongrpc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

type ForceBuyInRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ForceBuyInRequest) Reset() {
	*x = ForceBuyInRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForceBuyInRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForceBuyInRequest) ProtoMessage() {}

func (x *ForceBuyInRequest) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForceBuyInRequest.ProtoReflect.Descriptor instead.
func (*ForceBuyInRequest) Descriptor() ([]byte, []int) {
	return file_txpoker_action_service_proto_rawDescGZIP(), []int{0}
}

type ForceBuyInResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsBuyIn    bool                 `protobuf:"varint,1,opt,name=is_buy_in,json=isBuyIn,proto3" json:"is_buy_in,omitempty"`
	BuyInChip  int32                `protobuf:"varint,2,opt,name=buy_in_chip,json=buyInChip,proto3" json:"buy_in_chip,omitempty"`
	RemainTime *durationpb.Duration `protobuf:"bytes,3,opt,name=remainTime,proto3" json:"remainTime,omitempty"`
}

func (x *ForceBuyInResponse) Reset() {
	*x = ForceBuyInResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForceBuyInResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForceBuyInResponse) ProtoMessage() {}

func (x *ForceBuyInResponse) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForceBuyInResponse.ProtoReflect.Descriptor instead.
func (*ForceBuyInResponse) Descriptor() ([]byte, []int) {
	return file_txpoker_action_service_proto_rawDescGZIP(), []int{1}
}

func (x *ForceBuyInResponse) GetIsBuyIn() bool {
	if x != nil {
		return x.IsBuyIn
	}
	return false
}

func (x *ForceBuyInResponse) GetBuyInChip() int32 {
	if x != nil {
		return x.BuyInChip
	}
	return 0
}

func (x *ForceBuyInResponse) GetRemainTime() *durationpb.Duration {
	if x != nil {
		return x.RemainTime
	}
	return nil
}

var File_txpoker_action_service_proto protoreflect.FileDescriptor

var file_txpoker_action_service_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f,
	0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x73, 0x79, 0x6e, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x13, 0x0a, 0x11,
	0x46, 0x6f, 0x72, 0x63, 0x65, 0x42, 0x75, 0x79, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x8b, 0x01, 0x0a, 0x12, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x42, 0x75, 0x79, 0x49, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x62,
	0x75, 0x79, 0x5f, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x42,
	0x75, 0x79, 0x49, 0x6e, 0x12, 0x1e, 0x0a, 0x0b, 0x62, 0x75, 0x79, 0x5f, 0x69, 0x6e, 0x5f, 0x63,
	0x68, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x62, 0x75, 0x79, 0x49, 0x6e,
	0x43, 0x68, 0x69, 0x70, 0x12, 0x39, 0x0a, 0x0a, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x72, 0x65, 0x6d, 0x61, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x32,
	0xf9, 0x09, 0x0a, 0x0d, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x39, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x79, 0x6e, 0x63, 0x12, 0x15, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x37, 0x0a, 0x05,
	0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52,
	0x65, 0x61, 0x64, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3f, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x47, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x04, 0x4b, 0x69, 0x63, 0x6b, 0x12, 0x13,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4b, 0x69, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x41, 0x0a,
	0x0a, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x19, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x6f, 0x6f, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x3c, 0x0a, 0x07, 0x53, 0x74, 0x61, 0x6e, 0x64, 0x55, 0x70, 0x12, 0x17, 0x2e, 0x74, 0x78,
	0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x6e, 0x64, 0x55, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c,
	0x0a, 0x07, 0x53, 0x69, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x12, 0x17, 0x2e, 0x74, 0x78, 0x70, 0x6f,
	0x6b, 0x65, 0x72, 0x2e, 0x53, 0x69, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x05,
	0x42, 0x75, 0x79, 0x49, 0x6e, 0x12, 0x15, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x42, 0x75, 0x79, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x06, 0x53, 0x69, 0x74, 0x4f, 0x75, 0x74,
	0x12, 0x16, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x69, 0x74, 0x4f, 0x75,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x38, 0x0a, 0x05, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x12, 0x15, 0x2e, 0x74, 0x78,
	0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x04,
	0x46, 0x6f, 0x6c, 0x64, 0x12, 0x14, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x46,
	0x6f, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x05, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x15, 0x2e,
	0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x34,
	0x0a, 0x03, 0x42, 0x65, 0x74, 0x12, 0x13, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x42, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x04, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x14, 0x2e, 0x74,
	0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x05,
	0x52, 0x61, 0x69, 0x73, 0x65, 0x12, 0x15, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x52, 0x61, 0x69, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x05, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x12,
	0x15, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x3e, 0x0a, 0x08, 0x53, 0x68, 0x6f, 0x77, 0x46, 0x6f, 0x6c, 0x64, 0x12, 0x18, 0x2e, 0x74,
	0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x68, 0x6f, 0x77, 0x46, 0x6f, 0x6c, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00,
	0x12, 0x54, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x61, 0x69, 0x74, 0x42, 0x42,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x23, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x57, 0x61, 0x69, 0x74, 0x42, 0x42, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x41, 0x75, 0x74, 0x6f, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x12, 0x26, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x41, 0x75, 0x74, 0x6f, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x47, 0x0a, 0x0a, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x42, 0x75, 0x79, 0x49, 0x6e,
	0x12, 0x1a, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x46, 0x6f, 0x72, 0x63, 0x65,
	0x42, 0x75, 0x79, 0x49, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x74,
	0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x42, 0x75, 0x79, 0x49,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x6d, 0x5a, 0x2f, 0x63,
	0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2f, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02,
	0x1f, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e,
	0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0xca, 0xb2, 0x04, 0x16, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c,
	0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_txpoker_action_service_proto_rawDescOnce sync.Once
	file_txpoker_action_service_proto_rawDescData = file_txpoker_action_service_proto_rawDesc
)

func file_txpoker_action_service_proto_rawDescGZIP() []byte {
	file_txpoker_action_service_proto_rawDescOnce.Do(func() {
		file_txpoker_action_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_txpoker_action_service_proto_rawDescData)
	})
	return file_txpoker_action_service_proto_rawDescData
}

var file_txpoker_action_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_txpoker_action_service_proto_goTypes = []interface{}{
	(*ForceBuyInRequest)(nil),             // 0: txpoker.ForceBuyInRequest
	(*ForceBuyInResponse)(nil),            // 1: txpoker.ForceBuyInResponse
	(*durationpb.Duration)(nil),           // 2: google.protobuf.Duration
	(*commongrpc.ResyncRequest)(nil),      // 3: common.ResyncRequest
	(*commongrpc.ReadyRequest)(nil),       // 4: common.ReadyRequest
	(*commongrpc.StartGameRequest)(nil),   // 5: common.StartGameRequest
	(*commongrpc.KickRequest)(nil),        // 6: common.KickRequest
	(*commongrpc.ChangeRoomRequest)(nil),  // 7: common.ChangeRoomRequest
	(*StandUpRequest)(nil),                // 8: txpoker.StandUpRequest
	(*SitDownRequest)(nil),                // 9: txpoker.SitDownRequest
	(*BuyInRequest)(nil),                  // 10: txpoker.BuyInRequest
	(*SitOutRequest)(nil),                 // 11: txpoker.SitOutRequest
	(*TopUpRequest)(nil),                  // 12: txpoker.TopUpRequest
	(*FoldRequest)(nil),                   // 13: txpoker.FoldRequest
	(*CheckRequest)(nil),                  // 14: txpoker.CheckRequest
	(*BetRequest)(nil),                    // 15: txpoker.BetRequest
	(*CallRequest)(nil),                   // 16: txpoker.CallRequest
	(*RaiseRequest)(nil),                  // 17: txpoker.RaiseRequest
	(*AllInRequest)(nil),                  // 18: txpoker.AllInRequest
	(*ShowFoldRequest)(nil),               // 19: txpoker.ShowFoldRequest
	(*UpdateWaitBBSettingRequest)(nil),    // 20: txpoker.UpdateWaitBBSettingRequest
	(*UpdateAutoTopUpSettingRequest)(nil), // 21: txpoker.UpdateAutoTopUpSettingRequest
	(*emptypb.Empty)(nil),                 // 22: google.protobuf.Empty
}
var file_txpoker_action_service_proto_depIdxs = []int32{
	2,  // 0: txpoker.ForceBuyInResponse.remainTime:type_name -> google.protobuf.Duration
	3,  // 1: txpoker.ActionService.Resync:input_type -> common.ResyncRequest
	4,  // 2: txpoker.ActionService.Ready:input_type -> common.ReadyRequest
	5,  // 3: txpoker.ActionService.StartGame:input_type -> common.StartGameRequest
	6,  // 4: txpoker.ActionService.Kick:input_type -> common.KickRequest
	7,  // 5: txpoker.ActionService.ChangeRoom:input_type -> common.ChangeRoomRequest
	8,  // 6: txpoker.ActionService.StandUp:input_type -> txpoker.StandUpRequest
	9,  // 7: txpoker.ActionService.SitDown:input_type -> txpoker.SitDownRequest
	10, // 8: txpoker.ActionService.BuyIn:input_type -> txpoker.BuyInRequest
	11, // 9: txpoker.ActionService.SitOut:input_type -> txpoker.SitOutRequest
	12, // 10: txpoker.ActionService.TopUp:input_type -> txpoker.TopUpRequest
	13, // 11: txpoker.ActionService.Fold:input_type -> txpoker.FoldRequest
	14, // 12: txpoker.ActionService.Check:input_type -> txpoker.CheckRequest
	15, // 13: txpoker.ActionService.Bet:input_type -> txpoker.BetRequest
	16, // 14: txpoker.ActionService.Call:input_type -> txpoker.CallRequest
	17, // 15: txpoker.ActionService.Raise:input_type -> txpoker.RaiseRequest
	18, // 16: txpoker.ActionService.AllIn:input_type -> txpoker.AllInRequest
	19, // 17: txpoker.ActionService.ShowFold:input_type -> txpoker.ShowFoldRequest
	20, // 18: txpoker.ActionService.UpdateWaitBBSetting:input_type -> txpoker.UpdateWaitBBSettingRequest
	21, // 19: txpoker.ActionService.UpdateAutoTopUpSetting:input_type -> txpoker.UpdateAutoTopUpSettingRequest
	0,  // 20: txpoker.ActionService.ForceBuyIn:input_type -> txpoker.ForceBuyInRequest
	22, // 21: txpoker.ActionService.Resync:output_type -> google.protobuf.Empty
	22, // 22: txpoker.ActionService.Ready:output_type -> google.protobuf.Empty
	22, // 23: txpoker.ActionService.StartGame:output_type -> google.protobuf.Empty
	22, // 24: txpoker.ActionService.Kick:output_type -> google.protobuf.Empty
	22, // 25: txpoker.ActionService.ChangeRoom:output_type -> google.protobuf.Empty
	22, // 26: txpoker.ActionService.StandUp:output_type -> google.protobuf.Empty
	22, // 27: txpoker.ActionService.SitDown:output_type -> google.protobuf.Empty
	22, // 28: txpoker.ActionService.BuyIn:output_type -> google.protobuf.Empty
	22, // 29: txpoker.ActionService.SitOut:output_type -> google.protobuf.Empty
	22, // 30: txpoker.ActionService.TopUp:output_type -> google.protobuf.Empty
	22, // 31: txpoker.ActionService.Fold:output_type -> google.protobuf.Empty
	22, // 32: txpoker.ActionService.Check:output_type -> google.protobuf.Empty
	22, // 33: txpoker.ActionService.Bet:output_type -> google.protobuf.Empty
	22, // 34: txpoker.ActionService.Call:output_type -> google.protobuf.Empty
	22, // 35: txpoker.ActionService.Raise:output_type -> google.protobuf.Empty
	22, // 36: txpoker.ActionService.AllIn:output_type -> google.protobuf.Empty
	22, // 37: txpoker.ActionService.ShowFold:output_type -> google.protobuf.Empty
	22, // 38: txpoker.ActionService.UpdateWaitBBSetting:output_type -> google.protobuf.Empty
	22, // 39: txpoker.ActionService.UpdateAutoTopUpSetting:output_type -> google.protobuf.Empty
	1,  // 40: txpoker.ActionService.ForceBuyIn:output_type -> txpoker.ForceBuyInResponse
	21, // [21:41] is the sub-list for method output_type
	1,  // [1:21] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_txpoker_action_service_proto_init() }
func file_txpoker_action_service_proto_init() {
	if File_txpoker_action_service_proto != nil {
		return
	}
	file_txpoker_action_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_txpoker_action_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForceBuyInRequest); i {
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
		file_txpoker_action_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForceBuyInResponse); i {
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
			RawDescriptor: file_txpoker_action_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_txpoker_action_service_proto_goTypes,
		DependencyIndexes: file_txpoker_action_service_proto_depIdxs,
		MessageInfos:      file_txpoker_action_service_proto_msgTypes,
	}.Build()
	File_txpoker_action_service_proto = out.File
	file_txpoker_action_service_proto_rawDesc = nil
	file_txpoker_action_service_proto_goTypes = nil
	file_txpoker_action_service_proto_depIdxs = nil
}
