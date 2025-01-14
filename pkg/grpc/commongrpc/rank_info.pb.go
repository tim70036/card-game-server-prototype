// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: common/rank_info.proto

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

type RankInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid               string `protobuf:"bytes,1,opt,name=Uid,proto3" json:"Uid,omitempty"`
	BeforeRank        int32  `protobuf:"varint,2,opt,name=BeforeRank,proto3" json:"BeforeRank,omitempty"`
	AfterRank         int32  `protobuf:"varint,3,opt,name=AfterRank,proto3" json:"AfterRank,omitempty"`
	BeforeRankPoint   int32  `protobuf:"varint,4,opt,name=BeforeRankPoint,proto3" json:"BeforeRankPoint,omitempty"`
	PointChanged      int32  `protobuf:"varint,5,opt,name=PointChanged,proto3" json:"PointChanged,omitempty"`
	ExtraPointChanged int32  `protobuf:"varint,6,opt,name=ExtraPointChanged,proto3" json:"ExtraPointChanged,omitempty"`
	AfterRankPoint    int32  `protobuf:"varint,7,opt,name=AfterRankPoint,proto3" json:"AfterRankPoint,omitempty"`
	BeforeRate        int32  `protobuf:"varint,8,opt,name=BeforeRate,proto3" json:"BeforeRate,omitempty"`
	RateChanged       int32  `protobuf:"varint,9,opt,name=RateChanged,proto3" json:"RateChanged,omitempty"`
	BeforeRankUpPoint int32  `protobuf:"varint,10,opt,name=BeforeRankUpPoint,proto3" json:"BeforeRankUpPoint,omitempty"`
	AfterRankUpPoint  int32  `protobuf:"varint,11,opt,name=AfterRankUpPoint,proto3" json:"AfterRankUpPoint,omitempty"`
}

func (x *RankInfo) Reset() {
	*x = RankInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_rank_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RankInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RankInfo) ProtoMessage() {}

func (x *RankInfo) ProtoReflect() protoreflect.Message {
	mi := &file_common_rank_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RankInfo.ProtoReflect.Descriptor instead.
func (*RankInfo) Descriptor() ([]byte, []int) {
	return file_common_rank_info_proto_rawDescGZIP(), []int{0}
}

func (x *RankInfo) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *RankInfo) GetBeforeRank() int32 {
	if x != nil {
		return x.BeforeRank
	}
	return 0
}

func (x *RankInfo) GetAfterRank() int32 {
	if x != nil {
		return x.AfterRank
	}
	return 0
}

func (x *RankInfo) GetBeforeRankPoint() int32 {
	if x != nil {
		return x.BeforeRankPoint
	}
	return 0
}

func (x *RankInfo) GetPointChanged() int32 {
	if x != nil {
		return x.PointChanged
	}
	return 0
}

func (x *RankInfo) GetExtraPointChanged() int32 {
	if x != nil {
		return x.ExtraPointChanged
	}
	return 0
}

func (x *RankInfo) GetAfterRankPoint() int32 {
	if x != nil {
		return x.AfterRankPoint
	}
	return 0
}

func (x *RankInfo) GetBeforeRate() int32 {
	if x != nil {
		return x.BeforeRate
	}
	return 0
}

func (x *RankInfo) GetRateChanged() int32 {
	if x != nil {
		return x.RateChanged
	}
	return 0
}

func (x *RankInfo) GetBeforeRankUpPoint() int32 {
	if x != nil {
		return x.BeforeRankUpPoint
	}
	return 0
}

func (x *RankInfo) GetAfterRankUpPoint() int32 {
	if x != nil {
		return x.AfterRankUpPoint
	}
	return 0
}

var File_common_rank_info_proto protoreflect.FileDescriptor

var file_common_rank_info_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x72, 0x61, 0x6e, 0x6b, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x03, 0x0a, 0x08, 0x52, 0x61, 0x6e, 0x6b, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x55, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x52, 0x61,
	0x6e, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65,
	0x52, 0x61, 0x6e, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x66, 0x74, 0x65, 0x72, 0x52, 0x61, 0x6e,
	0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x41, 0x66, 0x74, 0x65, 0x72, 0x52, 0x61,
	0x6e, 0x6b, 0x12, 0x28, 0x0a, 0x0f, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x52, 0x61, 0x6e, 0x6b,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x42, 0x65, 0x66,
	0x6f, 0x72, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64,
	0x12, 0x2c, 0x0a, 0x11, 0x45, 0x78, 0x74, 0x72, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x45, 0x78, 0x74,
	0x72, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x26,
	0x0a, 0x0e, 0x41, 0x66, 0x74, 0x65, 0x72, 0x52, 0x61, 0x6e, 0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x41, 0x66, 0x74, 0x65, 0x72, 0x52, 0x61, 0x6e,
	0x6b, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65,
	0x52, 0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x42, 0x65, 0x66, 0x6f,
	0x72, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x52, 0x61, 0x74, 0x65, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x52, 0x61, 0x74,
	0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x64, 0x12, 0x2c, 0x0a, 0x11, 0x42, 0x65, 0x66, 0x6f,
	0x72, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x55, 0x70, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x11, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x52, 0x61, 0x6e, 0x6b, 0x55,
	0x70, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x2a, 0x0a, 0x10, 0x41, 0x66, 0x74, 0x65, 0x72, 0x52,
	0x61, 0x6e, 0x6b, 0x55, 0x70, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x10, 0x41, 0x66, 0x74, 0x65, 0x72, 0x52, 0x61, 0x6e, 0x6b, 0x55, 0x70, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x42, 0x6a, 0x5a, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1a, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d,
	0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0xca, 0xb2, 0x04, 0x15, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70,
	0x6c, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0xd0, 0xb2, 0x04, 0x01, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_rank_info_proto_rawDescOnce sync.Once
	file_common_rank_info_proto_rawDescData = file_common_rank_info_proto_rawDesc
)

func file_common_rank_info_proto_rawDescGZIP() []byte {
	file_common_rank_info_proto_rawDescOnce.Do(func() {
		file_common_rank_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_rank_info_proto_rawDescData)
	})
	return file_common_rank_info_proto_rawDescData
}

var file_common_rank_info_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_common_rank_info_proto_goTypes = []interface{}{
	(*RankInfo)(nil), // 0: common.RankInfo
}
var file_common_rank_info_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_rank_info_proto_init() }
func file_common_rank_info_proto_init() {
	if File_common_rank_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_rank_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RankInfo); i {
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
			RawDescriptor: file_common_rank_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_rank_info_proto_goTypes,
		DependencyIndexes: file_common_rank_info_proto_depIdxs,
		MessageInfos:      file_common_rank_info_proto_msgTypes,
	}.Build()
	File_common_rank_info_proto = out.File
	file_common_rank_info_proto_rawDesc = nil
	file_common_rank_info_proto_goTypes = nil
	file_common_rank_info_proto_depIdxs = nil
}
