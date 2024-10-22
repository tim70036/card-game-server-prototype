package service

import (
	"card-game-server-prototype/pkg/core"

	"google.golang.org/protobuf/proto"
)

// todo: 這裡有個疑問。
//  先 Init 再 FetchFromRepo，如果 fail 了，需要移除已經 init 的資料 or LeaveRoom 嗎？？
//  1. 照理說應該必要 api 打完沒問題後再將資料統整拿去 init user，但改動太大，先記著。
//  2. Destroy 是 for user 更新 and BroadcastUpdate 後使用。（裡面有 Unicast），
//     缺少還沒 BroadcastUpdate，只是單純移除 init 資料不需要 Unicast 的功能。
//     像是 Add AI 的時候，完全成功才會 Broadcast
//  3. BroadcastUpdate 做法也怪怪的，像是 buddy room 要廣播 2 次。

type UserService interface {
	Init(uids ...core.Uid) error
	GetAI(enterLimit, roomCardLimit int) (core.Uid, error)
	FetchFromRepo(uids ...core.Uid) error
	Destroy(topic core.Topic, kickoutMSG proto.Message, uids ...core.Uid) error
	CanConnect(uid core.Uid) error
	BroadcastUpdate()
}
