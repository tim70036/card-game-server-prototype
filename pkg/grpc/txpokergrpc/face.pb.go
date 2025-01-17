// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: txpoker/face.proto

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

type Face int32

const (
	Face_UNDEFINED Face = 0
	Face_ACE       Face = 1
	Face_TWO       Face = 2
	Face_THREE     Face = 3
	Face_FOUR      Face = 4
	Face_FIVE      Face = 5
	Face_SIX       Face = 6
	Face_SEVEN     Face = 7
	Face_EIGHT     Face = 8
	Face_NINE      Face = 9
	Face_TEN       Face = 10
	Face_JACK      Face = 11
	Face_QUEEN     Face = 12
	Face_KING      Face = 13
)

// Enum value maps for Face.
var (
	Face_name = map[int32]string{
		0:  "UNDEFINED",
		1:  "ACE",
		2:  "TWO",
		3:  "THREE",
		4:  "FOUR",
		5:  "FIVE",
		6:  "SIX",
		7:  "SEVEN",
		8:  "EIGHT",
		9:  "NINE",
		10: "TEN",
		11: "JACK",
		12: "QUEEN",
		13: "KING",
	}
	Face_value = map[string]int32{
		"UNDEFINED": 0,
		"ACE":       1,
		"TWO":       2,
		"THREE":     3,
		"FOUR":      4,
		"FIVE":      5,
		"SIX":       6,
		"SEVEN":     7,
		"EIGHT":     8,
		"NINE":      9,
		"TEN":       10,
		"JACK":      11,
		"QUEEN":     12,
		"KING":      13,
	}
)

func (x Face) Enum() *Face {
	p := new(Face)
	*p = x
	return p
}

func (x Face) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Face) Descriptor() protoreflect.EnumDescriptor {
	return file_txpoker_face_proto_enumTypes[0].Descriptor()
}

func (Face) Type() protoreflect.EnumType {
	return &file_txpoker_face_proto_enumTypes[0]
}

func (x Face) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Face.Descriptor instead.
func (Face) EnumDescriptor() ([]byte, []int) {
	return file_txpoker_face_proto_rawDescGZIP(), []int{0}
}

var File_txpoker_face_proto protoreflect.FileDescriptor

var file_txpoker_face_proto_rawDesc = []byte{
	0x0a, 0x12, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x66, 0x61,
	0x63, 0x65, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x97, 0x01, 0x0a, 0x04, 0x46, 0x61, 0x63, 0x65,
	0x12, 0x0d, 0x0a, 0x09, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x07, 0x0a, 0x03, 0x41, 0x43, 0x45, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x57, 0x4f, 0x10,
	0x02, 0x12, 0x09, 0x0a, 0x05, 0x54, 0x48, 0x52, 0x45, 0x45, 0x10, 0x03, 0x12, 0x08, 0x0a, 0x04,
	0x46, 0x4f, 0x55, 0x52, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x49, 0x56, 0x45, 0x10, 0x05,
	0x12, 0x07, 0x0a, 0x03, 0x53, 0x49, 0x58, 0x10, 0x06, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x45, 0x56,
	0x45, 0x4e, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x49, 0x47, 0x48, 0x54, 0x10, 0x08, 0x12,
	0x08, 0x0a, 0x04, 0x4e, 0x49, 0x4e, 0x45, 0x10, 0x09, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x45, 0x4e,
	0x10, 0x0a, 0x12, 0x08, 0x0a, 0x04, 0x4a, 0x41, 0x43, 0x4b, 0x10, 0x0b, 0x12, 0x09, 0x0a, 0x05,
	0x51, 0x55, 0x45, 0x45, 0x4e, 0x10, 0x0c, 0x12, 0x08, 0x0a, 0x04, 0x4b, 0x49, 0x4e, 0x47, 0x10,
	0x0d, 0x42, 0x6b, 0x5a, 0x2f, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72,
	0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1d, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d,
	0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x4d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0xca, 0xb2, 0x04, 0x16, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61,
	0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_txpoker_face_proto_rawDescOnce sync.Once
	file_txpoker_face_proto_rawDescData = file_txpoker_face_proto_rawDesc
)

func file_txpoker_face_proto_rawDescGZIP() []byte {
	file_txpoker_face_proto_rawDescOnce.Do(func() {
		file_txpoker_face_proto_rawDescData = protoimpl.X.CompressGZIP(file_txpoker_face_proto_rawDescData)
	})
	return file_txpoker_face_proto_rawDescData
}

var file_txpoker_face_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_txpoker_face_proto_goTypes = []interface{}{
	(Face)(0), // 0: txpoker.face.Face
}
var file_txpoker_face_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_txpoker_face_proto_init() }
func file_txpoker_face_proto_init() {
	if File_txpoker_face_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_txpoker_face_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_txpoker_face_proto_goTypes,
		DependencyIndexes: file_txpoker_face_proto_depIdxs,
		EnumInfos:         file_txpoker_face_proto_enumTypes,
	}.Build()
	File_txpoker_face_proto = out.File
	file_txpoker_face_proto_rawDesc = nil
	file_txpoker_face_proto_goTypes = nil
	file_txpoker_face_proto_depIdxs = nil
}
