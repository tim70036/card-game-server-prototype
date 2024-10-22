// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: txpoker/poker_hand_type.proto

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

type PokerHandType int32

const (
	PokerHandType_UNDEFINED       PokerHandType = 0
	PokerHandType_HIGH_CARD       PokerHandType = 1
	PokerHandType_PAIR            PokerHandType = 2
	PokerHandType_TWO_PAIR        PokerHandType = 3
	PokerHandType_THREE_OF_A_KIND PokerHandType = 4
	PokerHandType_STRAIGHT        PokerHandType = 5
	PokerHandType_FLUSH           PokerHandType = 6
	PokerHandType_FULL_HOUSE      PokerHandType = 7
	PokerHandType_FOUR_OF_A_KIND  PokerHandType = 8
	PokerHandType_STRAIGHT_FLUSH  PokerHandType = 9
	PokerHandType_ROYAL_FLUSH     PokerHandType = 10
)

// Enum value maps for PokerHandType.
var (
	PokerHandType_name = map[int32]string{
		0:  "UNDEFINED",
		1:  "HIGH_CARD",
		2:  "PAIR",
		3:  "TWO_PAIR",
		4:  "THREE_OF_A_KIND",
		5:  "STRAIGHT",
		6:  "FLUSH",
		7:  "FULL_HOUSE",
		8:  "FOUR_OF_A_KIND",
		9:  "STRAIGHT_FLUSH",
		10: "ROYAL_FLUSH",
	}
	PokerHandType_value = map[string]int32{
		"UNDEFINED":       0,
		"HIGH_CARD":       1,
		"PAIR":            2,
		"TWO_PAIR":        3,
		"THREE_OF_A_KIND": 4,
		"STRAIGHT":        5,
		"FLUSH":           6,
		"FULL_HOUSE":      7,
		"FOUR_OF_A_KIND":  8,
		"STRAIGHT_FLUSH":  9,
		"ROYAL_FLUSH":     10,
	}
)

func (x PokerHandType) Enum() *PokerHandType {
	p := new(PokerHandType)
	*p = x
	return p
}

func (x PokerHandType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PokerHandType) Descriptor() protoreflect.EnumDescriptor {
	return file_txpoker_poker_hand_type_proto_enumTypes[0].Descriptor()
}

func (PokerHandType) Type() protoreflect.EnumType {
	return &file_txpoker_poker_hand_type_proto_enumTypes[0]
}

func (x PokerHandType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PokerHandType.Descriptor instead.
func (PokerHandType) EnumDescriptor() ([]byte, []int) {
	return file_txpoker_poker_hand_type_proto_rawDescGZIP(), []int{0}
}

var File_txpoker_poker_hand_type_proto protoreflect.FileDescriptor

var file_txpoker_poker_hand_type_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x5f,
	0x68, 0x61, 0x6e, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x15, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x68, 0x61,
	0x6e, 0x64, 0x74, 0x79, 0x70, 0x65, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69,
	0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xbc, 0x01, 0x0a, 0x0d,
	0x50, 0x6f, 0x6b, 0x65, 0x72, 0x48, 0x61, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a,
	0x09, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09,
	0x48, 0x49, 0x47, 0x48, 0x5f, 0x43, 0x41, 0x52, 0x44, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x50,
	0x41, 0x49, 0x52, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x57, 0x4f, 0x5f, 0x50, 0x41, 0x49,
	0x52, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x54, 0x48, 0x52, 0x45, 0x45, 0x5f, 0x4f, 0x46, 0x5f,
	0x41, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x54, 0x52, 0x41,
	0x49, 0x47, 0x48, 0x54, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x4c, 0x55, 0x53, 0x48, 0x10,
	0x06, 0x12, 0x0e, 0x0a, 0x0a, 0x46, 0x55, 0x4c, 0x4c, 0x5f, 0x48, 0x4f, 0x55, 0x53, 0x45, 0x10,
	0x07, 0x12, 0x12, 0x0a, 0x0e, 0x46, 0x4f, 0x55, 0x52, 0x5f, 0x4f, 0x46, 0x5f, 0x41, 0x5f, 0x4b,
	0x49, 0x4e, 0x44, 0x10, 0x08, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x54, 0x52, 0x41, 0x49, 0x47, 0x48,
	0x54, 0x5f, 0x46, 0x4c, 0x55, 0x53, 0x48, 0x10, 0x09, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x4f, 0x59,
	0x41, 0x4c, 0x5f, 0x46, 0x4c, 0x55, 0x53, 0x48, 0x10, 0x0a, 0x42, 0x64, 0x5a, 0x2f, 0x63, 0x61,
	0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x16,
	0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54,
	0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0xca, 0xb2, 0x04, 0x16, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_txpoker_poker_hand_type_proto_rawDescOnce sync.Once
	file_txpoker_poker_hand_type_proto_rawDescData = file_txpoker_poker_hand_type_proto_rawDesc
)

func file_txpoker_poker_hand_type_proto_rawDescGZIP() []byte {
	file_txpoker_poker_hand_type_proto_rawDescOnce.Do(func() {
		file_txpoker_poker_hand_type_proto_rawDescData = protoimpl.X.CompressGZIP(file_txpoker_poker_hand_type_proto_rawDescData)
	})
	return file_txpoker_poker_hand_type_proto_rawDescData
}

var file_txpoker_poker_hand_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_txpoker_poker_hand_type_proto_goTypes = []interface{}{
	(PokerHandType)(0), // 0: txpoker.pokerhandtype.PokerHandType
}
var file_txpoker_poker_hand_type_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_txpoker_poker_hand_type_proto_init() }
func file_txpoker_poker_hand_type_proto_init() {
	if File_txpoker_poker_hand_type_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_txpoker_poker_hand_type_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_txpoker_poker_hand_type_proto_goTypes,
		DependencyIndexes: file_txpoker_poker_hand_type_proto_depIdxs,
		EnumInfos:         file_txpoker_poker_hand_type_proto_enumTypes,
	}.Build()
	File_txpoker_poker_hand_type_proto = out.File
	file_txpoker_poker_hand_type_proto_rawDesc = nil
	file_txpoker_poker_hand_type_proto_goTypes = nil
	file_txpoker_poker_hand_type_proto_depIdxs = nil
}
