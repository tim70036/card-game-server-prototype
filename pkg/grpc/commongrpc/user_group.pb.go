// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: common/user_group.proto

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

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid         string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Username    string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	IsConnected bool   `protobuf:"varint,3,opt,name=is_connected,json=isConnected,proto3" json:"is_connected,omitempty"`
	HasEntered  bool   `protobuf:"varint,4,opt,name=has_entered,json=hasEntered,proto3" json:"has_entered,omitempty"`
	Cash        int32  `protobuf:"varint,5,opt,name=cash,proto3" json:"cash,omitempty"`
	Level       int32  `protobuf:"varint,6,opt,name=level,proto3" json:"level,omitempty"`
	RoomCards   int32  `protobuf:"varint,7,opt,name=room_cards,json=roomCards,proto3" json:"room_cards,omitempty"`
	ShortUid    string `protobuf:"bytes,8,opt,name=short_uid,json=shortUid,proto3" json:"short_uid,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_user_group_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_common_user_group_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_common_user_group_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetIsConnected() bool {
	if x != nil {
		return x.IsConnected
	}
	return false
}

func (x *User) GetHasEntered() bool {
	if x != nil {
		return x.HasEntered
	}
	return false
}

func (x *User) GetCash() int32 {
	if x != nil {
		return x.Cash
	}
	return 0
}

func (x *User) GetLevel() int32 {
	if x != nil {
		return x.Level
	}
	return 0
}

func (x *User) GetRoomCards() int32 {
	if x != nil {
		return x.RoomCards
	}
	return 0
}

func (x *User) GetShortUid() string {
	if x != nil {
		return x.ShortUid
	}
	return ""
}

type UserGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users map[string]*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *UserGroup) Reset() {
	*x = UserGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_user_group_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserGroup) ProtoMessage() {}

func (x *UserGroup) ProtoReflect() protoreflect.Message {
	mi := &file_common_user_group_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserGroup.ProtoReflect.Descriptor instead.
func (*UserGroup) Descriptor() ([]byte, []int) {
	return file_common_user_group_proto_rawDescGZIP(), []int{1}
}

func (x *UserGroup) GetUsers() map[string]*User {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_common_user_group_proto protoreflect.FileDescriptor

var file_common_user_group_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xde, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x69, 0x73, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x61, 0x73, 0x5f, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x65, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x68, 0x61, 0x73, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x65,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x63, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x72,
	0x6f, 0x6f, 0x6d, 0x5f, 0x63, 0x61, 0x72, 0x64, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x72, 0x6f, 0x6f, 0x6d, 0x43, 0x61, 0x72, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x55, 0x69, 0x64, 0x22, 0x87, 0x01, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x32, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x46, 0x0a, 0x0a, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x22, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x42, 0x6a, 0x5a, 0x2e, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x67,
	0x72, 0x70, 0x63, 0xaa, 0x02, 0x1a, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65,
	0x70, 0x6c, 0x61, 0x79, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0xca, 0xb2, 0x04, 0x15, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c,
	0x61, 0x79, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0xd0, 0xb2, 0x04, 0x01, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_user_group_proto_rawDescOnce sync.Once
	file_common_user_group_proto_rawDescData = file_common_user_group_proto_rawDesc
)

func file_common_user_group_proto_rawDescGZIP() []byte {
	file_common_user_group_proto_rawDescOnce.Do(func() {
		file_common_user_group_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_user_group_proto_rawDescData)
	})
	return file_common_user_group_proto_rawDescData
}

var file_common_user_group_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_common_user_group_proto_goTypes = []interface{}{
	(*User)(nil),      // 0: common.User
	(*UserGroup)(nil), // 1: common.UserGroup
	nil,               // 2: common.UserGroup.UsersEntry
}
var file_common_user_group_proto_depIdxs = []int32{
	2, // 0: common.UserGroup.users:type_name -> common.UserGroup.UsersEntry
	0, // 1: common.UserGroup.UsersEntry.value:type_name -> common.User
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_common_user_group_proto_init() }
func file_common_user_group_proto_init() {
	if File_common_user_group_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_user_group_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_common_user_group_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserGroup); i {
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
			RawDescriptor: file_common_user_group_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_user_group_proto_goTypes,
		DependencyIndexes: file_common_user_group_proto_depIdxs,
		MessageInfos:      file_common_user_group_proto_msgTypes,
	}.Build()
	File_common_user_group_proto = out.File
	file_common_user_group_proto_rawDesc = nil
	file_common_user_group_proto_goTypes = nil
	file_common_user_group_proto_depIdxs = nil
}
