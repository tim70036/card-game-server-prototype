// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: txpoker/action_hint_group.proto

package txpokergrpc

import (
	_ "card-game-server-prototype/pkg/grpc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ActionHint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid              string               `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	BetChip          int32                `protobuf:"varint,2,opt,name=bet_chip,json=betChip,proto3" json:"bet_chip,omitempty"`
	CallingChip      int32                `protobuf:"varint,3,opt,name=calling_chip,json=callingChip,proto3" json:"calling_chip,omitempty"`
	MinRaiseChip     int32                `protobuf:"varint,4,opt,name=min_raise_chip,json=minRaiseChip,proto3" json:"min_raise_chip,omitempty"`
	Action           BetActionType        `protobuf:"varint,5,opt,name=action,proto3,enum=txpoker.betActionType.BetActionType" json:"action,omitempty"`
	AvailableActions []BetActionType      `protobuf:"varint,6,rep,packed,name=available_actions,json=availableActions,proto3,enum=txpoker.betActionType.BetActionType" json:"available_actions,omitempty"`
	Duration         *durationpb.Duration `protobuf:"bytes,7,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *ActionHint) Reset() {
	*x = ActionHint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_hint_group_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionHint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionHint) ProtoMessage() {}

func (x *ActionHint) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_hint_group_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionHint.ProtoReflect.Descriptor instead.
func (*ActionHint) Descriptor() ([]byte, []int) {
	return file_txpoker_action_hint_group_proto_rawDescGZIP(), []int{0}
}

func (x *ActionHint) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *ActionHint) GetBetChip() int32 {
	if x != nil {
		return x.BetChip
	}
	return 0
}

func (x *ActionHint) GetCallingChip() int32 {
	if x != nil {
		return x.CallingChip
	}
	return 0
}

func (x *ActionHint) GetMinRaiseChip() int32 {
	if x != nil {
		return x.MinRaiseChip
	}
	return 0
}

func (x *ActionHint) GetAction() BetActionType {
	if x != nil {
		return x.Action
	}
	return BetActionType_UNDEFINED
}

func (x *ActionHint) GetAvailableActions() []BetActionType {
	if x != nil {
		return x.AvailableActions
	}
	return nil
}

func (x *ActionHint) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

type ActionHintGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hints      map[string]*ActionHint `protobuf:"bytes,1,rep,name=hints,proto3" json:"hints,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RaiserHint *ActionHint            `protobuf:"bytes,2,opt,name=raiser_hint,json=raiserHint,proto3" json:"raiser_hint,omitempty"`
}

func (x *ActionHintGroup) Reset() {
	*x = ActionHintGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_txpoker_action_hint_group_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionHintGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionHintGroup) ProtoMessage() {}

func (x *ActionHintGroup) ProtoReflect() protoreflect.Message {
	mi := &file_txpoker_action_hint_group_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionHintGroup.ProtoReflect.Descriptor instead.
func (*ActionHintGroup) Descriptor() ([]byte, []int) {
	return file_txpoker_action_hint_group_proto_rawDescGZIP(), []int{1}
}

func (x *ActionHintGroup) GetHints() map[string]*ActionHint {
	if x != nil {
		return x.Hints
	}
	return nil
}

func (x *ActionHintGroup) GetRaiserHint() *ActionHint {
	if x != nil {
		return x.RaiserHint
	}
	return nil
}

var File_txpoker_action_hint_group_proto protoreflect.FileDescriptor

var file_txpoker_action_hint_group_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x68, 0x69, 0x6e, 0x74, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x07, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62,
	0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1d, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2f, 0x62, 0x65, 0x74, 0x5f, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xca,
	0x02, 0x0a, 0x0a, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x69, 0x6e, 0x74, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12,
	0x19, 0x0a, 0x08, 0x62, 0x65, 0x74, 0x5f, 0x63, 0x68, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x62, 0x65, 0x74, 0x43, 0x68, 0x69, 0x70, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61,
	0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x68, 0x69, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0b, 0x63, 0x61, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x69, 0x70, 0x12, 0x24, 0x0a,
	0x0e, 0x6d, 0x69, 0x6e, 0x5f, 0x72, 0x61, 0x69, 0x73, 0x65, 0x5f, 0x63, 0x68, 0x69, 0x70, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x6d, 0x69, 0x6e, 0x52, 0x61, 0x69, 0x73, 0x65, 0x43,
	0x68, 0x69, 0x70, 0x12, 0x3c, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x62, 0x65,
	0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x2e, 0x42, 0x65, 0x74, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x51, 0x0a, 0x11, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x24, 0x2e, 0x74,
	0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x62, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x2e, 0x42, 0x65, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x10, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xd1, 0x01, 0x0a, 0x0f,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x69, 0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x39, 0x0a, 0x05, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x69, 0x6e, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x48, 0x69, 0x6e, 0x74, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x05, 0x68, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x34, 0x0a, 0x0b, 0x72, 0x61,
	0x69, 0x73, 0x65, 0x72, 0x5f, 0x68, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x48, 0x69, 0x6e, 0x74, 0x52, 0x0a, 0x72, 0x61, 0x69, 0x73, 0x65, 0x72, 0x48, 0x69, 0x6e, 0x74,
	0x1a, 0x4d, 0x0a, 0x0a, 0x48, 0x69, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x29, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x48, 0x69, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42,
	0x6b, 0x5a, 0x2f, 0x63, 0x61, 0x72, 0x64, 0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x78, 0x70, 0x6f, 0x6b, 0x65, 0x72, 0x67, 0x72,
	0x70, 0x63, 0xaa, 0x02, 0x1d, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70,
	0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x4d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0xca, 0xb2, 0x04, 0x16, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65,
	0x70, 0x6c, 0x61, 0x79, 0x2e, 0x54, 0x78, 0x50, 0x6f, 0x6b, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_txpoker_action_hint_group_proto_rawDescOnce sync.Once
	file_txpoker_action_hint_group_proto_rawDescData = file_txpoker_action_hint_group_proto_rawDesc
)

func file_txpoker_action_hint_group_proto_rawDescGZIP() []byte {
	file_txpoker_action_hint_group_proto_rawDescOnce.Do(func() {
		file_txpoker_action_hint_group_proto_rawDescData = protoimpl.X.CompressGZIP(file_txpoker_action_hint_group_proto_rawDescData)
	})
	return file_txpoker_action_hint_group_proto_rawDescData
}

var file_txpoker_action_hint_group_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_txpoker_action_hint_group_proto_goTypes = []interface{}{
	(*ActionHint)(nil),          // 0: txpoker.ActionHint
	(*ActionHintGroup)(nil),     // 1: txpoker.ActionHintGroup
	nil,                         // 2: txpoker.ActionHintGroup.HintsEntry
	(BetActionType)(0),          // 3: txpoker.betActionType.BetActionType
	(*durationpb.Duration)(nil), // 4: google.protobuf.Duration
}
var file_txpoker_action_hint_group_proto_depIdxs = []int32{
	3, // 0: txpoker.ActionHint.action:type_name -> txpoker.betActionType.BetActionType
	3, // 1: txpoker.ActionHint.available_actions:type_name -> txpoker.betActionType.BetActionType
	4, // 2: txpoker.ActionHint.duration:type_name -> google.protobuf.Duration
	2, // 3: txpoker.ActionHintGroup.hints:type_name -> txpoker.ActionHintGroup.HintsEntry
	0, // 4: txpoker.ActionHintGroup.raiser_hint:type_name -> txpoker.ActionHint
	0, // 5: txpoker.ActionHintGroup.HintsEntry.value:type_name -> txpoker.ActionHint
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_txpoker_action_hint_group_proto_init() }
func file_txpoker_action_hint_group_proto_init() {
	if File_txpoker_action_hint_group_proto != nil {
		return
	}
	file_txpoker_bet_action_type_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_txpoker_action_hint_group_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActionHint); i {
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
		file_txpoker_action_hint_group_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActionHintGroup); i {
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
			RawDescriptor: file_txpoker_action_hint_group_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_txpoker_action_hint_group_proto_goTypes,
		DependencyIndexes: file_txpoker_action_hint_group_proto_depIdxs,
		MessageInfos:      file_txpoker_action_hint_group_proto_msgTypes,
	}.Build()
	File_txpoker_action_hint_group_proto = out.File
	file_txpoker_action_hint_group_proto_rawDesc = nil
	file_txpoker_action_hint_group_proto_goTypes = nil
	file_txpoker_action_hint_group_proto_depIdxs = nil
}
