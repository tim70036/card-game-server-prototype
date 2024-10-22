// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: common/room_info.proto

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

type RoomInfo_GameType int32

const (
	RoomInfo_CMJ           RoomInfo_GameType = 0 // todo: CMJ 搬移後，這個要改成 undefined/-1
	RoomInfo_DMJ           RoomInfo_GameType = 1
	RoomInfo_YABLON        RoomInfo_GameType = 2
	RoomInfo_CN_POKER      RoomInfo_GameType = 3
	RoomInfo_MJ            RoomInfo_GameType = 4
	RoomInfo_AMJ           RoomInfo_GameType = 5
	RoomInfo_TX_POKER      RoomInfo_GameType = 6
	RoomInfo_DARK_CHESS    RoomInfo_GameType = 7
	RoomInfo_HORSE_RACE_MJ RoomInfo_GameType = 8
	RoomInfo_TX_POKER_ZOOM RoomInfo_GameType = 9
)

// Enum value maps for RoomInfo_GameType.
var (
	RoomInfo_GameType_name = map[int32]string{
		0: "CMJ",
		1: "DMJ",
		2: "YABLON",
		3: "CN_POKER",
		4: "MJ",
		5: "AMJ",
		6: "TX_POKER",
		7: "DARK_CHESS",
		8: "HORSE_RACE_MJ",
		9: "TX_POKER_ZOOM",
	}
	RoomInfo_GameType_value = map[string]int32{
		"CMJ":           0,
		"DMJ":           1,
		"YABLON":        2,
		"CN_POKER":      3,
		"MJ":            4,
		"AMJ":           5,
		"TX_POKER":      6,
		"DARK_CHESS":    7,
		"HORSE_RACE_MJ": 8,
		"TX_POKER_ZOOM": 9,
	}
)

func (x RoomInfo_GameType) Enum() *RoomInfo_GameType {
	p := new(RoomInfo_GameType)
	*p = x
	return p
}

func (x RoomInfo_GameType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RoomInfo_GameType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_room_info_proto_enumTypes[0].Descriptor()
}

func (RoomInfo_GameType) Type() protoreflect.EnumType {
	return &file_common_room_info_proto_enumTypes[0]
}

func (x RoomInfo_GameType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RoomInfo_GameType.Descriptor instead.
func (RoomInfo_GameType) EnumDescriptor() ([]byte, []int) {
	return file_common_room_info_proto_rawDescGZIP(), []int{0, 0}
}

type RoomInfo_GameMode int32

const (
	RoomInfo_BUDDY       RoomInfo_GameMode = 0
	RoomInfo_COMMON      RoomInfo_GameMode = 1
	RoomInfo_CLUB        RoomInfo_GameMode = 2
	RoomInfo_RANK        RoomInfo_GameMode = 3
	RoomInfo_CARNIVAL    RoomInfo_GameMode = 4
	RoomInfo_QUALIFIER   RoomInfo_GameMode = 5
	RoomInfo_ELIMINATION RoomInfo_GameMode = 6
)

// Enum value maps for RoomInfo_GameMode.
var (
	RoomInfo_GameMode_name = map[int32]string{
		0: "BUDDY",
		1: "COMMON",
		2: "CLUB",
		3: "RANK",
		4: "CARNIVAL",
		5: "QUALIFIER",
		6: "ELIMINATION",
	}
	RoomInfo_GameMode_value = map[string]int32{
		"BUDDY":       0,
		"COMMON":      1,
		"CLUB":        2,
		"RANK":        3,
		"CARNIVAL":    4,
		"QUALIFIER":   5,
		"ELIMINATION": 6,
	}
)

func (x RoomInfo_GameMode) Enum() *RoomInfo_GameMode {
	p := new(RoomInfo_GameMode)
	*p = x
	return p
}

func (x RoomInfo_GameMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RoomInfo_GameMode) Descriptor() protoreflect.EnumDescriptor {
	return file_common_room_info_proto_enumTypes[1].Descriptor()
}

func (RoomInfo_GameMode) Type() protoreflect.EnumType {
	return &file_common_room_info_proto_enumTypes[1]
}

func (x RoomInfo_GameMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RoomInfo_GameMode.Descriptor instead.
func (RoomInfo_GameMode) EnumDescriptor() ([]byte, []int) {
	return file_common_room_info_proto_rawDescGZIP(), []int{0, 1}
}

type RoomInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomId       string            `protobuf:"bytes,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
	ShortRoomId  string            `protobuf:"bytes,2,opt,name=short_room_id,json=shortRoomId,proto3" json:"short_room_id,omitempty"`
	RoomGameType RoomInfo_GameType `protobuf:"varint,3,opt,name=room_game_type,json=roomGameType,proto3,enum=common.RoomInfo_GameType" json:"room_game_type,omitempty"`
	RoomGameMode RoomInfo_GameMode `protobuf:"varint,4,opt,name=room_game_mode,json=roomGameMode,proto3,enum=common.RoomInfo_GameMode" json:"room_game_mode,omitempty"`
	GameMetaUid  string            `protobuf:"bytes,5,opt,name=game_meta_uid,json=gameMetaUid,proto3" json:"game_meta_uid,omitempty"`
	IsPremium    bool              `protobuf:"varint,6,opt,name=is_premium,json=isPremium,proto3" json:"is_premium,omitempty"`
}

func (x *RoomInfo) Reset() {
	*x = RoomInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_room_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomInfo) ProtoMessage() {}

func (x *RoomInfo) ProtoReflect() protoreflect.Message {
	mi := &file_common_room_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomInfo.ProtoReflect.Descriptor instead.
func (*RoomInfo) Descriptor() ([]byte, []int) {
	return file_common_room_info_proto_rawDescGZIP(), []int{0}
}

func (x *RoomInfo) GetRoomId() string {
	if x != nil {
		return x.RoomId
	}
	return ""
}

func (x *RoomInfo) GetShortRoomId() string {
	if x != nil {
		return x.ShortRoomId
	}
	return ""
}

func (x *RoomInfo) GetRoomGameType() RoomInfo_GameType {
	if x != nil {
		return x.RoomGameType
	}
	return RoomInfo_CMJ
}

func (x *RoomInfo) GetRoomGameMode() RoomInfo_GameMode {
	if x != nil {
		return x.RoomGameMode
	}
	return RoomInfo_BUDDY
}

func (x *RoomInfo) GetGameMetaUid() string {
	if x != nil {
		return x.GameMetaUid
	}
	return ""
}

func (x *RoomInfo) GetIsPremium() bool {
	if x != nil {
		return x.IsPremium
	}
	return false
}

var File_common_room_info_proto protoreflect.FileDescriptor

var file_common_room_info_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xff, 0x03, 0x0a, 0x08, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x5f, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x12,
	0x3f, 0x0a, 0x0e, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0c, 0x72, 0x6f, 0x6f, 0x6d, 0x47, 0x61, 0x6d, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x3f, 0x0a, 0x0e, 0x72, 0x6f, 0x6f, 0x6d, 0x5f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6d, 0x6f,
	0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x4d,
	0x6f, 0x64, 0x65, 0x52, 0x0c, 0x72, 0x6f, 0x6f, 0x6d, 0x47, 0x61, 0x6d, 0x65, 0x4d, 0x6f, 0x64,
	0x65, 0x12, 0x22, 0x0a, 0x0d, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x5f, 0x75,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x67, 0x61, 0x6d, 0x65, 0x4d, 0x65,
	0x74, 0x61, 0x55, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x5f, 0x70, 0x72, 0x65, 0x6d,
	0x69, 0x75, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x50, 0x72, 0x65,
	0x6d, 0x69, 0x75, 0x6d, 0x22, 0x8b, 0x01, 0x0a, 0x08, 0x47, 0x61, 0x6d, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x07, 0x0a, 0x03, 0x43, 0x4d, 0x4a, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x4d,
	0x4a, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x59, 0x41, 0x42, 0x4c, 0x4f, 0x4e, 0x10, 0x02, 0x12,
	0x0c, 0x0a, 0x08, 0x43, 0x4e, 0x5f, 0x50, 0x4f, 0x4b, 0x45, 0x52, 0x10, 0x03, 0x12, 0x06, 0x0a,
	0x02, 0x4d, 0x4a, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x4d, 0x4a, 0x10, 0x05, 0x12, 0x0c,
	0x0a, 0x08, 0x54, 0x58, 0x5f, 0x50, 0x4f, 0x4b, 0x45, 0x52, 0x10, 0x06, 0x12, 0x0e, 0x0a, 0x0a,
	0x44, 0x41, 0x52, 0x4b, 0x5f, 0x43, 0x48, 0x45, 0x53, 0x53, 0x10, 0x07, 0x12, 0x11, 0x0a, 0x0d,
	0x48, 0x4f, 0x52, 0x53, 0x45, 0x5f, 0x52, 0x41, 0x43, 0x45, 0x5f, 0x4d, 0x4a, 0x10, 0x08, 0x12,
	0x11, 0x0a, 0x0d, 0x54, 0x58, 0x5f, 0x50, 0x4f, 0x4b, 0x45, 0x52, 0x5f, 0x5a, 0x4f, 0x4f, 0x4d,
	0x10, 0x09, 0x22, 0x63, 0x0a, 0x08, 0x47, 0x61, 0x6d, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x09,
	0x0a, 0x05, 0x42, 0x55, 0x44, 0x44, 0x59, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x4f, 0x4d,
	0x4d, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x43, 0x4c, 0x55, 0x42, 0x10, 0x02, 0x12,
	0x08, 0x0a, 0x04, 0x52, 0x41, 0x4e, 0x4b, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x41, 0x52,
	0x4e, 0x49, 0x56, 0x41, 0x4c, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x51, 0x55, 0x41, 0x4c, 0x49,
	0x46, 0x49, 0x45, 0x52, 0x10, 0x05, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x4c, 0x49, 0x4d, 0x49, 0x4e,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x06, 0x42, 0x6a, 0x5a, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2d,
	0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1a, 0x4a, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x2e,
	0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0xca, 0xb2, 0x04, 0x15, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0xd0,
	0xb2, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_room_info_proto_rawDescOnce sync.Once
	file_common_room_info_proto_rawDescData = file_common_room_info_proto_rawDesc
)

func file_common_room_info_proto_rawDescGZIP() []byte {
	file_common_room_info_proto_rawDescOnce.Do(func() {
		file_common_room_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_room_info_proto_rawDescData)
	})
	return file_common_room_info_proto_rawDescData
}

var file_common_room_info_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_common_room_info_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_room_info_proto_goTypes = []interface{}{
	(RoomInfo_GameType)(0), // 0: common.RoomInfo.GameType
	(RoomInfo_GameMode)(0), // 1: common.RoomInfo.GameMode
	(*RoomInfo)(nil),       // 2: common.RoomInfo
}
var file_common_room_info_proto_depIdxs = []int32{
	0, // 0: common.RoomInfo.room_game_type:type_name -> common.RoomInfo.GameType
	1, // 1: common.RoomInfo.room_game_mode:type_name -> common.RoomInfo.GameMode
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_common_room_info_proto_init() }
func file_common_room_info_proto_init() {
	if File_common_room_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_room_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomInfo); i {
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
			RawDescriptor: file_common_room_info_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_room_info_proto_goTypes,
		DependencyIndexes: file_common_room_info_proto_depIdxs,
		EnumInfos:         file_common_room_info_proto_enumTypes,
		MessageInfos:      file_common_room_info_proto_msgTypes,
	}.Build()
	File_common_room_info_proto = out.File
	file_common_room_info_proto_rawDesc = nil
	file_common_room_info_proto_goTypes = nil
	file_common_room_info_proto_depIdxs = nil
}