// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: common/kickout.proto

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

type KickoutReason int32

const (
	KickoutReason_GAME_EXCEPTION       KickoutReason = 0
	KickoutReason_GAME_CLOSED          KickoutReason = 1
	KickoutReason_NOT_ENOUGH_CASH      KickoutReason = 2
	KickoutReason_USER_REQUESTED       KickoutReason = 3
	KickoutReason_OWNER_REQUESTED      KickoutReason = 4
	KickoutReason_CHANGE_ROOM          KickoutReason = 5
	KickoutReason_ELIMINATION_ABORT    KickoutReason = 6
	KickoutReason_MAINTENANCE          KickoutReason = 7
	KickoutReason_IDLE_TIMEOUT         KickoutReason = 8
	KickoutReason_ROOM_EXPIRED         KickoutReason = 9
	KickoutReason_NOT_ENOUGH_ROOM_CARD KickoutReason = 10
)

// Enum value maps for KickoutReason.
var (
	KickoutReason_name = map[int32]string{
		0:  "GAME_EXCEPTION",
		1:  "GAME_CLOSED",
		2:  "NOT_ENOUGH_CASH",
		3:  "USER_REQUESTED",
		4:  "OWNER_REQUESTED",
		5:  "CHANGE_ROOM",
		6:  "ELIMINATION_ABORT",
		7:  "MAINTENANCE",
		8:  "IDLE_TIMEOUT",
		9:  "ROOM_EXPIRED",
		10: "NOT_ENOUGH_ROOM_CARD",
	}
	KickoutReason_value = map[string]int32{
		"GAME_EXCEPTION":       0,
		"GAME_CLOSED":          1,
		"NOT_ENOUGH_CASH":      2,
		"USER_REQUESTED":       3,
		"OWNER_REQUESTED":      4,
		"CHANGE_ROOM":          5,
		"ELIMINATION_ABORT":    6,
		"MAINTENANCE":          7,
		"IDLE_TIMEOUT":         8,
		"ROOM_EXPIRED":         9,
		"NOT_ENOUGH_ROOM_CARD": 10,
	}
)

func (x KickoutReason) Enum() *KickoutReason {
	p := new(KickoutReason)
	*p = x
	return p
}

func (x KickoutReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (KickoutReason) Descriptor() protoreflect.EnumDescriptor {
	return file_common_kickout_proto_enumTypes[0].Descriptor()
}

func (KickoutReason) Type() protoreflect.EnumType {
	return &file_common_kickout_proto_enumTypes[0]
}

func (x KickoutReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use KickoutReason.Descriptor instead.
func (KickoutReason) EnumDescriptor() ([]byte, []int) {
	return file_common_kickout_proto_rawDescGZIP(), []int{0}
}

type Kickout struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reason KickoutReason `protobuf:"varint,1,opt,name=reason,proto3,enum=common.KickoutReason" json:"reason,omitempty"`
}

func (x *Kickout) Reset() {
	*x = Kickout{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_kickout_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Kickout) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Kickout) ProtoMessage() {}

func (x *Kickout) ProtoReflect() protoreflect.Message {
	mi := &file_common_kickout_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Kickout.ProtoReflect.Descriptor instead.
func (*Kickout) Descriptor() ([]byte, []int) {
	return file_common_kickout_proto_rawDescGZIP(), []int{0}
}

func (x *Kickout) GetReason() KickoutReason {
	if x != nil {
		return x.Reason
	}
	return KickoutReason_GAME_EXCEPTION
}

var File_common_kickout_proto protoreflect.FileDescriptor

var file_common_kickout_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6b, 0x69, 0x63, 0x6b, 0x6f, 0x75, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x1a, 0x13,
	0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x38, 0x0a, 0x07, 0x4b, 0x69, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x12, 0x2d,
	0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4b, 0x69, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2a, 0xe9, 0x01,
	0x0a, 0x0d, 0x4b, 0x69, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x0e, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x45, 0x58, 0x43, 0x45, 0x50, 0x54, 0x49, 0x4f,
	0x4e, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x43, 0x4c, 0x4f, 0x53,
	0x45, 0x44, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x4e, 0x4f, 0x55,
	0x47, 0x48, 0x5f, 0x43, 0x41, 0x53, 0x48, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x53, 0x45,
	0x52, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x45, 0x44, 0x10, 0x03, 0x12, 0x13, 0x0a,
	0x0f, 0x4f, 0x57, 0x4e, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x45, 0x44,
	0x10, 0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x5f, 0x52, 0x4f, 0x4f,
	0x4d, 0x10, 0x05, 0x12, 0x15, 0x0a, 0x11, 0x45, 0x4c, 0x49, 0x4d, 0x49, 0x4e, 0x41, 0x54, 0x49,
	0x4f, 0x4e, 0x5f, 0x41, 0x42, 0x4f, 0x52, 0x54, 0x10, 0x06, 0x12, 0x0f, 0x0a, 0x0b, 0x4d, 0x41,
	0x49, 0x4e, 0x54, 0x45, 0x4e, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x07, 0x12, 0x10, 0x0a, 0x0c, 0x49,
	0x44, 0x4c, 0x45, 0x5f, 0x54, 0x49, 0x4d, 0x45, 0x4f, 0x55, 0x54, 0x10, 0x08, 0x12, 0x10, 0x0a,
	0x0c, 0x52, 0x4f, 0x4f, 0x4d, 0x5f, 0x45, 0x58, 0x50, 0x49, 0x52, 0x45, 0x44, 0x10, 0x09, 0x12,
	0x18, 0x0a, 0x14, 0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x4e, 0x4f, 0x55, 0x47, 0x48, 0x5f, 0x52, 0x4f,
	0x4f, 0x4d, 0x5f, 0x43, 0x41, 0x52, 0x44, 0x10, 0x0a, 0x42, 0x6a, 0x5a, 0x2e, 0x63, 0x61, 0x72,
	0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1a, 0x4a, 0x6f,
	0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x47, 0x72, 0x70,
	0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0xca, 0xb2, 0x04, 0x15, 0x4a, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0xd0, 0xb2, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_kickout_proto_rawDescOnce sync.Once
	file_common_kickout_proto_rawDescData = file_common_kickout_proto_rawDesc
)

func file_common_kickout_proto_rawDescGZIP() []byte {
	file_common_kickout_proto_rawDescOnce.Do(func() {
		file_common_kickout_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_kickout_proto_rawDescData)
	})
	return file_common_kickout_proto_rawDescData
}

var file_common_kickout_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_kickout_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_kickout_proto_goTypes = []interface{}{
	(KickoutReason)(0), // 0: common.KickoutReason
	(*Kickout)(nil),    // 1: common.Kickout
}
var file_common_kickout_proto_depIdxs = []int32{
	0, // 0: common.Kickout.reason:type_name -> common.KickoutReason
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_common_kickout_proto_init() }
func file_common_kickout_proto_init() {
	if File_common_kickout_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_kickout_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Kickout); i {
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
			RawDescriptor: file_common_kickout_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_kickout_proto_goTypes,
		DependencyIndexes: file_common_kickout_proto_depIdxs,
		EnumInfos:         file_common_kickout_proto_enumTypes,
		MessageInfos:      file_common_kickout_proto_msgTypes,
	}.Build()
	File_common_kickout_proto = out.File
	file_common_kickout_proto_rawDesc = nil
	file_common_kickout_proto_goTypes = nil
	file_common_kickout_proto_depIdxs = nil
}
