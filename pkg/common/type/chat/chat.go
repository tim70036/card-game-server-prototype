package chat

import (
	commonconstant "card-game-server-prototype/pkg/common/constant"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"go.uber.org/zap/zapcore"
	"sync"
)

type Msg struct {
	Uid     core.Uid
	Message string
}

func NewMsg(uid core.Uid, m string) Msg {
	return Msg{
		Uid:     uid,
		Message: m,
	}
}

func (m Msg) ToProto() *commongrpc.Msg {
	return &commongrpc.Msg{
		Uid:     m.Uid.String(),
		Message: m.Message,
	}
}

func (m Msg) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", m.Uid.String())
	enc.AddString("message", m.Message)
	return nil
}

type History struct {
	mu   *sync.RWMutex
	data []Msg
}

func NewHistory() *History {
	return &History{
		mu:   &sync.RWMutex{},
		data: make([]Msg, 0),
	}
}

func (h *History) ToProto() []*commongrpc.Msg {
	var history []*commongrpc.Msg
	h.mu.RLock()
	for _, chat := range h.data {
		history = append(history, chat.ToProto())
	}
	h.mu.RUnlock()
	return history
}

func (h *History) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	h.mu.RLock()
	for _, v := range h.data {
		_ = enc.AppendObject(v)
	}
	h.mu.RUnlock()
	return nil
}

// todo: CD 時間

func (h *History) Add(m Msg) {
	// 暫時不阻擋重複
	// if h.isDuplicate(m) {
	// 	return
	// }

	h.mu.Lock()
	defer h.mu.Unlock()

	h.data = append([]Msg{m}, h.data...)

	if len(h.data) > commonconstant.ChatHistoryLimit {
		h.data = h.data[:commonconstant.ChatHistoryLimit]
	}
}

func (h *History) isDuplicate(m Msg) bool {
	if len(h.data) == 0 {
		return false
	}

	var last Msg
	h.mu.RLock()
	last = h.data[0]
	h.mu.RUnlock()

	return last.Uid == m.Uid && last.Message == m.Message
}
