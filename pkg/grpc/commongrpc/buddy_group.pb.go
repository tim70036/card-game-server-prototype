// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: common/buddy_group.proto

package commongrpc

import (
	_ "card-game-server-prototype/pkg/grpc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Buddy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid       string                 `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	IsReady   bool                   `protobuf:"varint,2,opt,name=is_ready,json=isReady,proto3" json:"is_ready,omitempty"`
	IsOwner   bool                   `protobuf:"varint,3,opt,name=is_owner,json=isOwner,proto3" json:"is_owner,omitempty"`
	EnterTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=enter_time,json=enterTime,proto3" json:"enter_time,omitempty"`
}

func (x *Buddy) Reset() {
	*x = Buddy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_buddy_group_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Buddy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Buddy) ProtoMessage() {}

func (x *Buddy) ProtoReflect() protoreflect.Message {
	mi := &file_common_buddy_group_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Buddy.ProtoReflect.Descriptor instead.
func (*Buddy) Descriptor() ([]byte, []int) {
	return file_common_buddy_group_proto_rawDescGZIP(), []int{0}
}

func (x *Buddy) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *Buddy) GetIsReady() bool {
	if x != nil {
		return x.IsReady
	}
	return false
}

func (x *Buddy) GetIsOwner() bool {
	if x != nil {
		return x.IsOwner
	}
	return false
}

func (x *Buddy) GetEnterTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EnterTime
	}
	return nil
}

type BuddyGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Buddies map[string]*Buddy `protobuf:"bytes,1,rep,name=buddies,proto3" json:"buddies,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *BuddyGroup) Reset() {
	*x = BuddyGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_buddy_group_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuddyGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuddyGroup) ProtoMessage() {}

func (x *BuddyGroup) ProtoReflect() protoreflect.Message {
	mi := &file_common_buddy_group_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuddyGroup.ProtoReflect.Descriptor instead.
func (*BuddyGroup) Descriptor() ([]byte, []int) {
	return file_common_buddy_group_proto_rawDescGZIP(), []int{1}
}

func (x *BuddyGroup) GetBuddies() map[string]*Buddy {
	if x != nil {
		return x.Buddies
	}
	return nil
}

var File_common_buddy_group_proto protoreflect.FileDescriptor

var file_common_buddy_group_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x62, 0x75, 0x64, 0x64, 0x79, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8a, 0x01, 0x0a, 0x05, 0x42, 0x75, 0x64, 0x64,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x19,
	0x0a, 0x08, 0x69, 0x73, 0x5f, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x69, 0x73, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x0a, 0x65, 0x6e, 0x74,
	0x65, 0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x65, 0x6e, 0x74, 0x65, 0x72,
	0x54, 0x69, 0x6d, 0x65, 0x22, 0x92, 0x01, 0x0a, 0x0a, 0x42, 0x75, 0x64, 0x64, 0x79, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x12, 0x39, 0x0a, 0x07, 0x62, 0x75, 0x64, 0x64, 0x69, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x42, 0x75,
	0x64, 0x64, 0x79, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x42, 0x75, 0x64, 0x64, 0x69, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x62, 0x75, 0x64, 0x64, 0x69, 0x65, 0x73, 0x1a, 0x49,
	0x0a, 0x0c, 0x42, 0x75, 0x64, 0x64, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x23, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x42, 0x75, 0x64, 0x64, 0x79, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x6a, 0x5a, 0x2e, 0x63, 0x61, 0x72,
	0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1a, 0x4a, 0x6f,
	0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x47, 0x72, 0x70,
	0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0xca, 0xb2, 0x04, 0x15, 0x4a, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0xd0, 0xb2, 0x04, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_buddy_group_proto_rawDescOnce sync.Once
	file_common_buddy_group_proto_rawDescData = file_common_buddy_group_proto_rawDesc
)

func file_common_buddy_group_proto_rawDescGZIP() []byte {
	file_common_buddy_group_proto_rawDescOnce.Do(func() {
		file_common_buddy_group_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_buddy_group_proto_rawDescData)
	})
	return file_common_buddy_group_proto_rawDescData
}

var file_common_buddy_group_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_common_buddy_group_proto_goTypes = []interface{}{
	(*Buddy)(nil),                 // 0: common.Buddy
	(*BuddyGroup)(nil),            // 1: common.BuddyGroup
	nil,                           // 2: common.BuddyGroup.BuddiesEntry
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_common_buddy_group_proto_depIdxs = []int32{
	3, // 0: common.Buddy.enter_time:type_name -> google.protobuf.Timestamp
	2, // 1: common.BuddyGroup.buddies:type_name -> common.BuddyGroup.BuddiesEntry
	0, // 2: common.BuddyGroup.BuddiesEntry.value:type_name -> common.Buddy
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_common_buddy_group_proto_init() }
func file_common_buddy_group_proto_init() {
	if File_common_buddy_group_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_buddy_group_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Buddy); i {
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
		file_common_buddy_group_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuddyGroup); i {
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
			RawDescriptor: file_common_buddy_group_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_buddy_group_proto_goTypes,
		DependencyIndexes: file_common_buddy_group_proto_depIdxs,
		MessageInfos:      file_common_buddy_group_proto_msgTypes,
	}.Build()
	File_common_buddy_group_proto = out.File
	file_common_buddy_group_proto_rawDesc = nil
	file_common_buddy_group_proto_goTypes = nil
	file_common_buddy_group_proto_depIdxs = nil
}
