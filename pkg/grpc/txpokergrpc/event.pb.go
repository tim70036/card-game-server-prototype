// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: txpoker/event.proto

package txpokergrpc

import (
	_ "card-game-server-prototype/pkg/grpc"
	commongrpc "card-game-server-prototype/pkg/grpc/commongrpc"
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

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kickout      *commongrpc.Kickout `protobuf:"bytes,1,opt,name=kickout,proto3,oneof" json:"kickout,omitempty"`
	ExpInfo      *commongrpc.ExpInfo `protobuf:"bytes,3,opt,name=exp_info,json=expInfo,proto3,oneof" json:"exp_info,omitempty"`
	Warning      *Warning            `protobuf:"bytes,4,opt,name=warning,proto3,oneof" json:"warning,omitempty"`
	TopupSuccess *TopUpSuccess       `protobuf:"bytes,5,opt,name=topup_success,json=topupSuccess,proto3,oneof" json:"topup_success,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_txpoker_event_proto_rawDescGZIP(), []int{0}
}

func (x *Event) GetKickout() *commongrpc.Kickout {
	if x != nil {
		return x.Kickout
	}
	return nil
}

func (x *Event) GetExpInfo() *commongrpc.ExpInfo {
	if x != nil {
		return x.ExpInfo
	}
	return nil
}

func (x *Event) GetWarning() *Warning {
	if x != nil {
		return x.Warning
	}
	return nil
}

func (x *Event) GetTopupSuccess() *TopUpSuccess {
	if x != nil {
		return x.TopupSuccess
	}
	return nil
}

type TopUpSuccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount int32 `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *TopUpSuccess) Reset() {
	*x = TopUpSuccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopUpSuccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopUpSuccess) ProtoMessage() {}

func (x *TopUpSuccess) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopUpSuccess.ProtoReflect.Descriptor instead.
func (*TopUpSuccess) Descriptor() ([]byte, []int) {
	return file_txpoker_event_proto_rawDescGZIP(), []int{1}
}

func (x *TopUpSuccess) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

var File_txpoker_event_proto protoreflect.FileDescriptor

var file_txpoker_event_proto_rawDesc = []byte{
	0x0a, 0x13, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x1a, 0x13,
	0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x6b, 0x69, 0x63, 0x6b,
	0x6f, 0x75, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x65, 0x78, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x15, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x12, 0x2e, 0x0a, 0x07, 0x6b, 0x69, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4b, 0x69, 0x63, 0x6b,
	0x6f, 0x75, 0x74, 0x48, 0x00, 0x52, 0x07, 0x6b, 0x69, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x88, 0x01,
	0x01, 0x12, 0x2f, 0x0a, 0x08, 0x65, 0x78, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x78, 0x70,
	0x49, 0x6e, 0x66, 0x6f, 0x48, 0x01, 0x52, 0x07, 0x65, 0x78, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x88,
	0x01, 0x01, 0x12, 0x2f, 0x0a, 0x07, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x57, 0x61,
	0x72, 0x6e, 0x69, 0x6e, 0x67, 0x48, 0x02, 0x52, 0x07, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67,
	0x88, 0x01, 0x01, 0x12, 0x3f, 0x0a, 0x0d, 0x74, 0x6f, 0x70, 0x75, 0x70, 0x5f, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x78, 0x70,
	0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x54, 0x6f, 0x70, 0x55, 0x70, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x48, 0x03, 0x52, 0x0c, 0x74, 0x6f, 0x70, 0x75, 0x70, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x6b, 0x69, 0x63, 0x6b, 0x6f, 0x75, 0x74,
	0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x65, 0x78, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x42, 0x0a, 0x0a,
	0x08, 0x5f, 0x77, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x74, 0x6f,
	0x70, 0x75, 0x70, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x26, 0x0a, 0x0c, 0x54,
	0x6f, 0x70, 0x55, 0x70, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x42, 0x6b, 0x5a, 0x2f, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65,
	0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70,
	0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x78, 0x70, 0x6f, 0x6b,
	0x65, 0x72, 0x67, 0x72, 0x70, 0x63, 0xaa, 0x02, 0x1d, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47,
	0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0xca, 0xb2, 0x04, 0x16, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_txpoker_event_proto_rawDescOnce sync.Once
	file_txpoker_event_proto_rawDescData = file_txpoker_event_proto_rawDesc
)

func file_txpoker_event_proto_rawDescGZIP() []byte {
	file_txpoker_event_proto_rawDescOnce.Do(func() {
		file_txpoker_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_txpoker_event_proto_rawDescData)
	})
	return file_txpoker_event_proto_rawDescData
}

var file_txpoker_event_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_txpoker_event_proto_goTypes = []interface{}{
	(*Event)(nil),              // 0: txpoker.Event
	(*TopUpSuccess)(nil),       // 1: txpoker.TopUpSuccess
	(*commongrpc.Kickout)(nil), // 2: common.Kickout
	(*commongrpc.ExpInfo)(nil), // 3: common.ExpInfo
	(*Warning)(nil),            // 4: txpoker.Warning
}
var file_txpoker_event_proto_depIdxs = []int32{
	2, // 0: txpoker.Event.kickout:type_name -> common.Kickout
	3, // 1: txpoker.Event.exp_info:type_name -> common.ExpInfo
	4, // 2: txpoker.Event.warning:type_name -> txpoker.Warning
	1, // 3: txpoker.Event.topup_success:type_name -> txpoker.TopUpSuccess
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_txpoker_event_proto_init() }
func file_txpoker_event_proto_init() {
	if File_txpoker_event_proto != nil {
		return
	}
	file_txpoker_warning_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_txpoker_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_txpoker_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopUpSuccess); i {
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
	file_txpoker_event_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_txpoker_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_txpoker_event_proto_goTypes,
		DependencyIndexes: file_txpoker_event_proto_depIdxs,
		MessageInfos:      file_txpoker_event_proto_msgTypes,
	}.Build()
	File_txpoker_event_proto = out.File
	file_txpoker_event_proto_rawDesc = nil
	file_txpoker_event_proto_goTypes = nil
	file_txpoker_event_proto_depIdxs = nil
}
