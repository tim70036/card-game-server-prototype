// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: extension/csharp_options.proto

package extensiongrpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CSharpAccess int32

const (
	CSharpAccess_internal CSharpAccess = 0
	CSharpAccess_public   CSharpAccess = 1
)

// Enum value maps for CSharpAccess.
var (
	CSharpAccess_name = map[int32]string{
		0: "internal",
		1: "public",
	}
	CSharpAccess_value = map[string]int32{
		"internal": 0,
		"public":   1,
	}
)

func (x CSharpAccess) Enum() *CSharpAccess {
	p := new(CSharpAccess)
	*p = x
	return p
}

func (x CSharpAccess) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CSharpAccess) Descriptor() protoreflect.EnumDescriptor {
	return file_extension_csharp_options_proto_enumTypes[0].Descriptor()
}

func (CSharpAccess) Type() protoreflect.EnumType {
	return &file_extension_csharp_options_proto_enumTypes[0]
}

func (x CSharpAccess) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CSharpAccess.Descriptor instead.
func (CSharpAccess) EnumDescriptor() ([]byte, []int) {
	return file_extension_csharp_options_proto_rawDescGZIP(), []int{0}
}

var file_extension_csharp_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         9001,
		Name:          "csharp_assembly",
		Tag:           "bytes,9001,opt,name=csharp_assembly",
		Filename:      "extension/csharp_options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*CSharpAccess)(nil),
		Field:         9002,
		Name:          "csharp_access",
		Tag:           "varint,9002,opt,name=csharp_access,enum=CSharpAccess",
		Filename:      "extension/csharp_options.proto",
	},
}

// Extension fields to descriptorpb.FileOptions.
var (
	// optional string csharp_assembly = 9001;
	E_CsharpAssembly = &file_extension_csharp_options_proto_extTypes[0]
	// optional CSharpAccess csharp_access = 9002;
	E_CsharpAccess = &file_extension_csharp_options_proto_extTypes[1]
)

var File_extension_csharp_options_proto protoreflect.FileDescriptor

var file_extension_csharp_options_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x63, 0x73, 0x68, 0x61,
	0x72, 0x70, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2a, 0x28, 0x0a, 0x0c, 0x43, 0x53, 0x68, 0x61, 0x72, 0x70, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x12, 0x0c, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x10, 0x00,
	0x12, 0x0a, 0x0a, 0x06, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x10, 0x01, 0x3a, 0x49, 0x0a, 0x0f,
	0x63, 0x73, 0x68, 0x61, 0x72, 0x70, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x6d, 0x62, 0x6c, 0x79, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xa9, 0x46,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x73, 0x68, 0x61, 0x72, 0x70, 0x41, 0x73, 0x73, 0x65,
	0x6d, 0x62, 0x6c, 0x79, 0x88, 0x01, 0x01, 0x3a, 0x54, 0x0a, 0x0d, 0x63, 0x73, 0x68, 0x61, 0x72,
	0x70, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xaa, 0x46, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e,
	0x43, 0x53, 0x68, 0x61, 0x72, 0x70, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x0c, 0x63, 0x73,
	0x68, 0x61, 0x72, 0x70, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x88, 0x01, 0x01, 0x42, 0x54, 0x5a,
	0x31, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x67, 0x72,
	0x70, 0x63, 0xaa, 0x02, 0x0a, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x72, 0x70, 0x63, 0xca,
	0xb2, 0x04, 0x0c, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0xd0,
	0xb2, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_extension_csharp_options_proto_rawDescOnce sync.Once
	file_extension_csharp_options_proto_rawDescData = file_extension_csharp_options_proto_rawDesc
)

func file_extension_csharp_options_proto_rawDescGZIP() []byte {
	file_extension_csharp_options_proto_rawDescOnce.Do(func() {
		file_extension_csharp_options_proto_rawDescData = protoimpl.X.CompressGZIP(file_extension_csharp_options_proto_rawDescData)
	})
	return file_extension_csharp_options_proto_rawDescData
}

var file_extension_csharp_options_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_extension_csharp_options_proto_goTypes = []interface{}{
	(CSharpAccess)(0),                // 0: CSharpAccess
	(*descriptorpb.FileOptions)(nil), // 1: google.protobuf.FileOptions
}
var file_extension_csharp_options_proto_depIdxs = []int32{
	1, // 0: csharp_assembly:extendee -> google.protobuf.FileOptions
	1, // 1: csharp_access:extendee -> google.protobuf.FileOptions
	0, // 2: csharp_access:type_name -> CSharpAccess
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	2, // [2:3] is the sub-list for extension type_name
	0, // [0:2] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_extension_csharp_options_proto_init() }
func file_extension_csharp_options_proto_init() {
	if File_extension_csharp_options_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_extension_csharp_options_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_extension_csharp_options_proto_goTypes,
		DependencyIndexes: file_extension_csharp_options_proto_depIdxs,
		EnumInfos:         file_extension_csharp_options_proto_enumTypes,
		ExtensionInfos:    file_extension_csharp_options_proto_extTypes,
	}.Build()
	File_extension_csharp_options_proto = out.File
	file_extension_csharp_options_proto_rawDesc = nil
	file_extension_csharp_options_proto_goTypes = nil
	file_extension_csharp_options_proto_depIdxs = nil
}