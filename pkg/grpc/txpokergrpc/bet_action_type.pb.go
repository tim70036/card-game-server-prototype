// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: txpoker/bet_action_type.proto

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

type BetActionType int32

const (
	BetActionType_UNDEFINED BetActionType = 0
	BetActionType_FOLD      BetActionType = 1
	BetActionType_CHECK     BetActionType = 2
	BetActionType_BET       BetActionType = 3
	BetActionType_CALL      BetActionType = 4
	BetActionType_RAISE     BetActionType = 5
	BetActionType_ALL_IN    BetActionType = 6
	BetActionType_SB        BetActionType = 7
	BetActionType_BB        BetActionType = 8
)

// Enum value maps for BetActionType.
var (
	BetActionType_name = map[int32]string{
		0: "UNDEFINED",
		1: "FOLD",
		2: "CHECK",
		3: "BET",
		4: "CALL",
		5: "RAISE",
		6: "ALL_IN",
		7: "SB",
		8: "BB",
	}
	BetActionType_value = map[string]int32{
		"UNDEFINED": 0,
		"FOLD":      1,
		"CHECK":     2,
		"BET":       3,
		"CALL":      4,
		"RAISE":     5,
		"ALL_IN":    6,
		"SB":        7,
		"BB":        8,
	}
)

func (x BetActionType) Enum() *BetActionType {
	p := new(BetActionType)
	*p = x
	return p
}

func (x BetActionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BetActionType) Descriptor() protoreflect.EnumDescriptor {
	return file_txpoker_bet_action_type_proto_enumTypes[0].Descriptor()
}

func (BetActionType) Type() protoreflect.EnumType {
	return &file_txpoker_bet_action_type_proto_enumTypes[0]
}

func (x BetActionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BetActionType.Descriptor instead.
func (BetActionType) EnumDescriptor() ([]byte, []int) {
	return file_txpoker_bet_action_type_proto_rawDescGZIP(), []int{0}
}

var File_txpoker_bet_action_type_proto protoreflect.FileDescriptor

var file_txpoker_bet_action_type_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x62, 0x65, 0x74, 0x5f, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x15, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x62, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x6d, 0x0a, 0x0d, 0x42,
	0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09,
	0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x46,
	0x4f, 0x4c, 0x44, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x48, 0x45, 0x43, 0x4b, 0x10, 0x02,
	0x12, 0x07, 0x0a, 0x03, 0x42, 0x45, 0x54, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04, 0x43, 0x41, 0x4c,
	0x4c, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x41, 0x49, 0x53, 0x45, 0x10, 0x05, 0x12, 0x0a,
	0x0a, 0x06, 0x41, 0x4c, 0x4c, 0x5f, 0x49, 0x4e, 0x10, 0x06, 0x12, 0x06, 0x0a, 0x02, 0x53, 0x42,
	0x10, 0x07, 0x12, 0x06, 0x0a, 0x02, 0x42, 0x42, 0x10, 0x08, 0x42, 0x64, 0x5a, 0x2f, 0x63, 0x61,
	0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x16,
	0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54,
	0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0xca, 0xb2, 0x04, 0x16, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_txpoker_bet_action_type_proto_rawDescOnce sync.Once
	file_txpoker_bet_action_type_proto_rawDescData = file_txpoker_bet_action_type_proto_rawDesc
)

func file_txpoker_bet_action_type_proto_rawDescGZIP() []byte {
	file_txpoker_bet_action_type_proto_rawDescOnce.Do(func() {
		file_txpoker_bet_action_type_proto_rawDescData = protoimpl.X.CompressGZIP(file_txpoker_bet_action_type_proto_rawDescData)
	})
	return file_txpoker_bet_action_type_proto_rawDescData
}

var file_txpoker_bet_action_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_txpoker_bet_action_type_proto_goTypes = []interface{}{
	(BetActionType)(0), // 0: txpoker.betActionType.BetActionType
}
var file_txpoker_bet_action_type_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_txpoker_bet_action_type_proto_init() }
func file_txpoker_bet_action_type_proto_init() {
	if File_txpoker_bet_action_type_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_txpoker_bet_action_type_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_txpoker_bet_action_type_proto_goTypes,
		DependencyIndexes: file_txpoker_bet_action_type_proto_depIdxs,
		EnumInfos:         file_txpoker_bet_action_type_proto_enumTypes,
	}.Build()
	File_txpoker_bet_action_type_proto = out.File
	file_txpoker_bet_action_type_proto_rawDesc = nil
	file_txpoker_bet_action_type_proto_goTypes = nil
	file_txpoker_bet_action_type_proto_depIdxs = nil
}
