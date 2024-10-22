package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/samber/lo"
)

type ResyncService struct {
	gameController core.GameController
	userGroup      *commonmodel.UserGroup
	buddyGroup     *commonmodel.BuddyGroup
	roomInfo       *commonmodel.RoomInfo

	gameInfo         *model2.GameInfo
	gameSetting      *model2.GameSetting
	table            *model2.Table
	seatStatusGroup  *model2.SeatStatusGroup
	actionHintGroup  *model2.ActionHintGroup
	playSettingGroup *model2.PlaySettingGroup
	statsCacheGroup  *model2.StatsCacheGroup
	chipCacheGroup   *model2.ChipCacheGroup
	userCacheGroup   *model2.UserCacheGroup
	playerGroup      *model2.PlayerGroup
	tableProfitGroup *model2.TableProfitsGroup

	msgBus core.MsgBus
}

func ProvideResyncService(
	gameController core.GameController,
	userGroup *commonmodel.UserGroup,
	buddyGroup *commonmodel.BuddyGroup,
	roomInfo *commonmodel.RoomInfo,

	gameInfo *model2.GameInfo,
	gameSetting *model2.GameSetting,
	table *model2.Table,
	seatStatusGroup *model2.SeatStatusGroup,
	actionHintGroup *model2.ActionHintGroup,
	playSettingGroup *model2.PlaySettingGroup,
	statsCacheGroup *model2.StatsCacheGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	userCacheGroup *model2.UserCacheGroup,
	playerGroup *model2.PlayerGroup,
	tableProfitGroup *model2.TableProfitsGroup,
	msgBus core.MsgBus,
) *ResyncService {
	return &ResyncService{
		gameController: gameController,
		userGroup:      userGroup,
		buddyGroup:     buddyGroup,
		roomInfo:       roomInfo,

		gameInfo:         gameInfo,
		gameSetting:      gameSetting,
		table:            table,
		seatStatusGroup:  seatStatusGroup,
		actionHintGroup:  actionHintGroup,
		playSettingGroup: playSettingGroup,
		statsCacheGroup:  statsCacheGroup,
		chipCacheGroup:   chipCacheGroup,
		userCacheGroup:   userCacheGroup,
		playerGroup:      playerGroup,
		tableProfitGroup: tableProfitGroup,
		msgBus:           msgBus,
	}
}

func (s *ResyncService) Send(uid core.Uid) {
	playerGroupProto := &txpokergrpc.PlayerGroup{Players: make(map[string]*txpokergrpc.Player)}
	for uid, player := range s.playerGroup.Data {
		playerGroupProto.Players[uid.String()] = player.ToProto()
		_, playerGroupProto.Players[uid.String()].HasShowdown = s.table.ShowdownPocketCards[uid]
	}

	modelMsg := &txpokergrpc.Model{
		UserGroup:       s.userGroup.ToProto(),
		BuddyGroup:      s.buddyGroup.ToProto(),
		RoomInfo:        s.roomInfo.ToProto(),
		ResyncGameState: s.gameController.CurrentState().ToProto(uid).(*txpokergrpc.GameState),

		GameInfo:          s.gameInfo.ToProto(),
		GameSetting:       s.gameSetting.ToProto(),
		Table:             s.table.ToProto(),
		SeatStatusGroup:   s.seatStatusGroup.ToProto(),
		PlayerGroup:       playerGroupProto,
		ActionHintGroup:   s.actionHintGroup.ToProto(),
		StatsCacheGroup:   s.statsCacheGroup.ToProto(),
		ChipCacheGroup:    s.chipCacheGroup.ToProto(),
		UserCacheGroup:    s.userCacheGroup.ToProto(),
		TableProfitsGroup: s.tableProfitGroup.ToProto(lo.Keys(s.playerGroup.Data)),
	}

	if playSetting, ok := s.playSettingGroup.Data[uid]; ok {
		modelMsg.PlaySetting = playSetting.ToProto()
	}

	s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
		Model: modelMsg,
	})
}
