// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: common/emote.proto

package commongrpc

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

type StickerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StickerId int32 `protobuf:"varint,1,opt,name=stickerId,proto3" json:"stickerId,omitempty"`
}

func (x *StickerRequest) Reset() {
	*x = StickerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_emote_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StickerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StickerRequest) ProtoMessage() {}

func (x *StickerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_emote_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StickerRequest.ProtoReflect.Descriptor instead.
func (*StickerRequest) Descriptor() ([]byte, []int) {
	return file_common_emote_proto_rawDescGZIP(), []int{0}
}

func (x *StickerRequest) GetStickerId() int32 {
	if x != nil {
		return x.StickerId
	}
	return 0
}

type Sticker struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid       string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	StickerId int32  `protobuf:"varint,2,opt,name=stickerId,proto3" json:"stickerId,omitempty"`
}

func (x *Sticker) Reset() {
	*x = Sticker{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_emote_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sticker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sticker) ProtoMessage() {}

func (x *Sticker) ProtoReflect() protoreflect.Message {
	mi := &file_common_emote_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sticker.ProtoReflect.Descriptor instead.
func (*Sticker) Descriptor() ([]byte, []int) {
	return file_common_emote_proto_rawDescGZIP(), []int{1}
}

func (x *Sticker) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *Sticker) GetStickerId() int32 {
	if x != nil {
		return x.StickerId
	}
	return 0
}

type EmotePingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemId    int32  `protobuf:"varint,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	TargetUid string `protobuf:"bytes,2,opt,name=target_uid,json=targetUid,proto3" json:"target_uid,omitempty"`
}

func (x *EmotePingRequest) Reset() {
	*x = EmotePingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_emote_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmotePingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmotePingRequest) ProtoMessage() {}

func (x *EmotePingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_common_emote_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmotePingRequest.ProtoReflect.Descriptor instead.
func (*EmotePingRequest) Descriptor() ([]byte, []int) {
	return file_common_emote_proto_rawDescGZIP(), []int{2}
}

func (x *EmotePingRequest) GetItemId() int32 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *EmotePingRequest) GetTargetUid() string {
	if x != nil {
		return x.TargetUid
	}
	return ""
}

type EmotePing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemId    int32  `protobuf:"varint,1,opt,name=itemId,proto3" json:"itemId,omitempty"`
	SenderUid string `protobuf:"bytes,2,opt,name=sender_uid,json=senderUid,proto3" json:"sender_uid,omitempty"`
	TargetUid string `protobuf:"bytes,3,opt,name=target_uid,json=targetUid,proto3" json:"target_uid,omitempty"`
}

func (x *EmotePing) Reset() {
	*x = EmotePing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_emote_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmotePing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmotePing) ProtoMessage() {}

func (x *EmotePing) ProtoReflect() protoreflect.Message {
	mi := &file_common_emote_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmotePing.ProtoReflect.Descriptor instead.
func (*EmotePing) Descriptor() ([]byte, []int) {
	return file_common_emote_proto_rawDescGZIP(), []int{3}
}

func (x *EmotePing) GetItemId() int32 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *EmotePing) GetSenderUid() string {
	if x != nil {
		return x.SenderUid
	}
	return ""
}

func (x *EmotePing) GetTargetUid() string {
	if x != nil {
		return x.TargetUid
	}
	return ""
}

type EmoteEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EmotePing *EmotePing `protobuf:"bytes,1,opt,name=emotePing,proto3,oneof" json:"emotePing,omitempty"`
	Sticker   *Sticker   `protobuf:"bytes,2,opt,name=sticker,proto3,oneof" json:"sticker,omitempty"`
}

func (x *EmoteEvent) Reset() {
	*x = EmoteEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_emote_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmoteEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmoteEvent) ProtoMessage() {}

func (x *EmoteEvent) ProtoReflect() protoreflect.Message {
	mi := &file_common_emote_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmoteEvent.ProtoReflect.Descriptor instead.
func (*EmoteEvent) Descriptor() ([]byte, []int) {
	return file_common_emote_proto_rawDescGZIP(), []int{4}
}

func (x *EmoteEvent) GetEmotePing() *EmotePing {
	if x != nil {
		return x.EmotePing
	}
	return nil
}

func (x *EmoteEvent) GetSticker() *Sticker {
	if x != nil {
		return x.Sticker
	}
	return nil
}

var File_common_emote_proto protoreflect.FileDescriptor

var file_common_emote_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x1a, 0x13, 0x67, 0x6c,
	0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x2e, 0x0a, 0x0e, 0x53, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x39, 0x0a, 0x07, 0x53, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x49, 0x64, 0x22, 0x49, 0x0a, 0x10,
	0x45, 0x6d, 0x6f, 0x74, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x55, 0x69, 0x64, 0x22, 0x61, 0x0a, 0x09, 0x45, 0x6d, 0x6f, 0x74, 0x65,
	0x50, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x69, 0x64, 0x22, 0x8c, 0x01, 0x0a, 0x0a, 0x45,
	0x6d, 0x6f, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x09, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x6d, 0x6f, 0x74, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x48,
	0x00, 0x52, 0x09, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12,
	0x2e, 0x0a, 0x07, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x69, 0x63, 0x6b, 0x65,
	0x72, 0x48, 0x01, 0x52, 0x07, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x88, 0x01, 0x01, 0x42,
	0x0c, 0x0a, 0x0a, 0x5f, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x73, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x42, 0x6a, 0x5a, 0x2e, 0x63, 0x61, 0x72,
	0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1a, 0x4a, 0x6f,
	0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x47, 0x72, 0x70,
	0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0xca, 0xb2, 0x04, 0x15, 0x4a, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0xd0, 0xb2, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_emote_proto_rawDescOnce sync.Once
	file_common_emote_proto_rawDescData = file_common_emote_proto_rawDesc
)

func file_common_emote_proto_rawDescGZIP() []byte {
	file_common_emote_proto_rawDescOnce.Do(func() {
		file_common_emote_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_emote_proto_rawDescData)
	})
	return file_common_emote_proto_rawDescData
}

var file_common_emote_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_common_emote_proto_goTypes = []interface{}{
	(*StickerRequest)(nil),   // 0: common.StickerRequest
	(*Sticker)(nil),          // 1: common.Sticker
	(*EmotePingRequest)(nil), // 2: common.EmotePingRequest
	(*EmotePing)(nil),        // 3: common.EmotePing
	(*EmoteEvent)(nil),       // 4: common.EmoteEvent
}
var file_common_emote_proto_depIdxs = []int32{
	3, // 0: common.EmoteEvent.emotePing:type_name -> common.EmotePing
	1, // 1: common.EmoteEvent.sticker:type_name -> common.Sticker
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_common_emote_proto_init() }
func file_common_emote_proto_init() {
	if File_common_emote_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_emote_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StickerRequest); i {
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
		file_common_emote_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sticker); i {
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
		file_common_emote_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmotePingRequest); i {
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
		file_common_emote_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmotePing); i {
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
		file_common_emote_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmoteEvent); i {
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
	file_common_emote_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_common_emote_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_emote_proto_goTypes,
		DependencyIndexes: file_common_emote_proto_depIdxs,
		MessageInfos:      file_common_emote_proto_msgTypes,
	}.Build()
	File_common_emote_proto = out.File
	file_common_emote_proto_rawDesc = nil
	file_common_emote_proto_goTypes = nil
	file_common_emote_proto_depIdxs = nil
}
