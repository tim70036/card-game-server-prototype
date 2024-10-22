// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: darkchess/scoreboard.proto

package darkchessgrpc

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

type RoundScore struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid            string            `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Points         int32             `protobuf:"varint,2,opt,name=points,proto3" json:"points,omitempty"`
	RawProfits     int32             `protobuf:"varint,3,opt,name=raw_profits,json=rawProfits,proto3" json:"raw_profits,omitempty"`
	CapturedPieces *CapturedPieces   `protobuf:"bytes,4,opt,name=captured_pieces,json=capturedPieces,proto3" json:"captured_pieces,omitempty"`
	ScoreModifier  ScoreModifierType `protobuf:"varint,5,opt,name=score_modifier,json=scoreModifier,proto3,enum=darkchess.ScoreModifierType" json:"score_modifier,omitempty"`
}

func (x *RoundScore) Reset() {
	*x = RoundScore{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_scoreboard_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoundScore) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoundScore) ProtoMessage() {}

func (x *RoundScore) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_scoreboard_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoundScore.ProtoReflect.Descriptor instead.
func (*RoundScore) Descriptor() ([]byte, []int) {
	return file_darkchess_scoreboard_proto_rawDescGZIP(), []int{0}
}

func (x *RoundScore) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *RoundScore) GetPoints() int32 {
	if x != nil {
		return x.Points
	}
	return 0
}

func (x *RoundScore) GetRawProfits() int32 {
	if x != nil {
		return x.RawProfits
	}
	return 0
}

func (x *RoundScore) GetCapturedPieces() *CapturedPieces {
	if x != nil {
		return x.CapturedPieces
	}
	return nil
}

func (x *RoundScore) GetScoreModifier() ScoreModifierType {
	if x != nil {
		return x.ScoreModifier
	}
	return ScoreModifierType_SCORE_MODIFIER_TYPE_INVALID
}

type RoundScoreboard struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scores    []*RoundScore `protobuf:"bytes,1,rep,name=scores,proto3" json:"scores,omitempty"`
	WinnerUid string        `protobuf:"bytes,2,opt,name=winner_uid,json=winnerUid,proto3" json:"winner_uid,omitempty"`
	IsDraw    bool          `protobuf:"varint,3,opt,name=is_draw,json=isDraw,proto3" json:"is_draw,omitempty"`
}

func (x *RoundScoreboard) Reset() {
	*x = RoundScoreboard{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_scoreboard_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoundScoreboard) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoundScoreboard) ProtoMessage() {}

func (x *RoundScoreboard) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_scoreboard_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoundScoreboard.ProtoReflect.Descriptor instead.
func (*RoundScoreboard) Descriptor() ([]byte, []int) {
	return file_darkchess_scoreboard_proto_rawDescGZIP(), []int{1}
}

func (x *RoundScoreboard) GetScores() []*RoundScore {
	if x != nil {
		return x.Scores
	}
	return nil
}

func (x *RoundScoreboard) GetWinnerUid() string {
	if x != nil {
		return x.WinnerUid
	}
	return ""
}

func (x *RoundScoreboard) GetIsDraw() bool {
	if x != nil {
		return x.IsDraw
	}
	return false
}

type GameScore struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid                string              `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Profit             int32               `protobuf:"varint,2,opt,name=profit,proto3" json:"profit,omitempty"`
	ExpInfo            *commongrpc.ExpInfo `protobuf:"bytes,3,opt,name=exp_info,json=expInfo,proto3" json:"exp_info,omitempty"`
	IsDisconnected     bool                `protobuf:"varint,4,opt,name=is_disconnected,json=isDisconnected,proto3" json:"is_disconnected,omitempty"`
	DisconnectedProfit int32               `protobuf:"varint,5,opt,name=disconnected_profit,json=disconnectedProfit,proto3" json:"disconnected_profit,omitempty"`
}

func (x *GameScore) Reset() {
	*x = GameScore{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_scoreboard_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameScore) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameScore) ProtoMessage() {}

func (x *GameScore) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_scoreboard_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameScore.ProtoReflect.Descriptor instead.
func (*GameScore) Descriptor() ([]byte, []int) {
	return file_darkchess_scoreboard_proto_rawDescGZIP(), []int{2}
}

func (x *GameScore) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *GameScore) GetProfit() int32 {
	if x != nil {
		return x.Profit
	}
	return 0
}

func (x *GameScore) GetExpInfo() *commongrpc.ExpInfo {
	if x != nil {
		return x.ExpInfo
	}
	return nil
}

func (x *GameScore) GetIsDisconnected() bool {
	if x != nil {
		return x.IsDisconnected
	}
	return false
}

func (x *GameScore) GetDisconnectedProfit() int32 {
	if x != nil {
		return x.DisconnectedProfit
	}
	return 0
}

type GameScoreboard struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scores    []*GameScore `protobuf:"bytes,1,rep,name=scores,proto3" json:"scores,omitempty"`
	WinnerUid string       `protobuf:"bytes,2,opt,name=winner_uid,json=winnerUid,proto3" json:"winner_uid,omitempty"`
	IsDraw    bool         `protobuf:"varint,3,opt,name=is_draw,json=isDraw,proto3" json:"is_draw,omitempty"`
}

func (x *GameScoreboard) Reset() {
	*x = GameScoreboard{}
	if protoimpl.UnsafeEnabled {
		mi := &file_darkchess_scoreboard_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameScoreboard) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameScoreboard) ProtoMessage() {}

func (x *GameScoreboard) ProtoReflect() protoreflect.Message {
	mi := &file_darkchess_scoreboard_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameScoreboard.ProtoReflect.Descriptor instead.
func (*GameScoreboard) Descriptor() ([]byte, []int) {
	return file_darkchess_scoreboard_proto_rawDescGZIP(), []int{3}
}

func (x *GameScoreboard) GetScores() []*GameScore {
	if x != nil {
		return x.Scores
	}
	return nil
}

func (x *GameScoreboard) GetWinnerUid() string {
	if x != nil {
		return x.WinnerUid
	}
	return ""
}

func (x *GameScoreboard) GetIsDraw() bool {
	if x != nil {
		return x.IsDraw
	}
	return false
}

var File_darkchess_scoreboard_proto protoreflect.FileDescriptor

var file_darkchess_scoreboard_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2f, 0x73, 0x63, 0x6f, 0x72,
	0x65, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x64, 0x61,
	0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x1a, 0x13, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f,
	0x69, 0x6d, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x65, 0x78, 0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2f, 0x65,
	0x6e, 0x75, 0x6d, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2f, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe0, 0x01, 0x0a, 0x0a, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x61, 0x77, 0x5f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x73, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x61, 0x77, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x73, 0x12,
	0x42, 0x0a, 0x0f, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x64, 0x5f, 0x70, 0x69, 0x65, 0x63,
	0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63,
	0x68, 0x65, 0x73, 0x73, 0x2e, 0x43, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x64, 0x50, 0x69, 0x65,
	0x63, 0x65, 0x73, 0x52, 0x0e, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x64, 0x50, 0x69, 0x65,
	0x63, 0x65, 0x73, 0x12, 0x43, 0x0a, 0x0e, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x6d, 0x6f, 0x64,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x64, 0x61,
	0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x4d, 0x6f, 0x64,
	0x69, 0x66, 0x69, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0d, 0x73, 0x63, 0x6f, 0x72, 0x65,
	0x4d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22, 0x78, 0x0a, 0x0f, 0x52, 0x6f, 0x75, 0x6e,
	0x64, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x62, 0x6f, 0x61, 0x72, 0x64, 0x12, 0x2d, 0x0a, 0x06, 0x73,
	0x63, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x64, 0x61,
	0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x53, 0x63, 0x6f,
	0x72, 0x65, 0x52, 0x06, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x77, 0x69,
	0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x55, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f,
	0x64, 0x72, 0x61, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x44, 0x72,
	0x61, 0x77, 0x22, 0xbb, 0x01, 0x0a, 0x09, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x63, 0x6f, 0x72, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x65, 0x78,
	0x70, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x45, 0x78, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x65,
	0x78, 0x70, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0e, 0x69, 0x73, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x12,
	0x2f, 0x0a, 0x13, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x70, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x12, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74,
	0x22, 0x76, 0x0a, 0x0e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x62, 0x6f, 0x61,
	0x72, 0x64, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x47,
	0x61, 0x6d, 0x65, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x06, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x73,
	0x12, 0x1d, 0x0a, 0x0a, 0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x75, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x77, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x55, 0x69, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x64, 0x72, 0x61, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x69, 0x73, 0x44, 0x72, 0x61, 0x77, 0x42, 0x72, 0x5a, 0x32, 0x63, 0x61, 0x72, 0x64,
	0x2d, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f,
	0x64, 0x61, 0x72, 0x6b, 0x63, 0x68, 0x65, 0x73, 0x73, 0x67, 0x72, 0x70, 0x63, 0x3b, 0xaa, 0x02,
	0x1f, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x2e,
	0x44, 0x61, 0x72, 0x6b, 0x43, 0x68, 0x65, 0x73, 0x73, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0xca, 0xb2, 0x04, 0x18, 0x4a, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x70, 0x6c,
	0x61, 0x79, 0x2e, 0x44, 0x61, 0x72, 0x6b, 0x43, 0x68, 0x65, 0x73, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_darkchess_scoreboard_proto_rawDescOnce sync.Once
	file_darkchess_scoreboard_proto_rawDescData = file_darkchess_scoreboard_proto_rawDesc
)

func file_darkchess_scoreboard_proto_rawDescGZIP() []byte {
	file_darkchess_scoreboard_proto_rawDescOnce.Do(func() {
		file_darkchess_scoreboard_proto_rawDescData = protoimpl.X.CompressGZIP(file_darkchess_scoreboard_proto_rawDescData)
	})
	return file_darkchess_scoreboard_proto_rawDescData
}

var file_darkchess_scoreboard_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_darkchess_scoreboard_proto_goTypes = []interface{}{
	(*RoundScore)(nil),         // 0: darkchess.RoundScore
	(*RoundScoreboard)(nil),    // 1: darkchess.RoundScoreboard
	(*GameScore)(nil),          // 2: darkchess.GameScore
	(*GameScoreboard)(nil),     // 3: darkchess.GameScoreboard
	(*CapturedPieces)(nil),     // 4: darkchess.CapturedPieces
	(ScoreModifierType)(0),     // 5: darkchess.ScoreModifierType
	(*commongrpc.ExpInfo)(nil), // 6: common.ExpInfo
}
var file_darkchess_scoreboard_proto_depIdxs = []int32{
	4, // 0: darkchess.RoundScore.captured_pieces:type_name -> darkchess.CapturedPieces
	5, // 1: darkchess.RoundScore.score_modifier:type_name -> darkchess.ScoreModifierType
	0, // 2: darkchess.RoundScoreboard.scores:type_name -> darkchess.RoundScore
	6, // 3: darkchess.GameScore.exp_info:type_name -> common.ExpInfo
	2, // 4: darkchess.GameScoreboard.scores:type_name -> darkchess.GameScore
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_darkchess_scoreboard_proto_init() }
func file_darkchess_scoreboard_proto_init() {
	if File_darkchess_scoreboard_proto != nil {
		return
	}
	file_darkchess_enum_enum_proto_init()
	file_darkchess_board_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_darkchess_scoreboard_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoundScore); i {
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
		file_darkchess_scoreboard_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoundScoreboard); i {
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
		file_darkchess_scoreboard_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameScore); i {
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
		file_darkchess_scoreboard_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameScoreboard); i {
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
			RawDescriptor: file_darkchess_scoreboard_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_darkchess_scoreboard_proto_goTypes,
		DependencyIndexes: file_darkchess_scoreboard_proto_depIdxs,
		MessageInfos:      file_darkchess_scoreboard_proto_msgTypes,
	}.Build()
	File_darkchess_scoreboard_proto = out.File
	file_darkchess_scoreboard_proto_rawDesc = nil
	file_darkchess_scoreboard_proto_goTypes = nil
	file_darkchess_scoreboard_proto_depIdxs = nil
}
