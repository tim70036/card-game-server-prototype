package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
)

type ResyncService struct {
	gameController   core.GameController
	userGroup        *commonmodel.UserGroup
	buddyGroup       *commonmodel.BuddyGroup
	playerGroup      *model2.PlayerGroup
	roomInfo         *commonmodel.RoomInfo
	playSettingGroup *model2.PlaySettingGroup
	gameInfo         *model2.GameInfo
	pickBoard        *model2.PickBoard
	board            *model2.Board
	capturedPieces   *model2.CapturedPieces
	actionHintGroup  *model2.ActionHintGroup

	msgBus core.MsgBus
}

func ProvideResyncService(
	gameController core.GameController,
	userGroup *commonmodel.UserGroup,
	buddyGroup *commonmodel.BuddyGroup,
	playerGroup *model2.PlayerGroup,
	roomInfo *commonmodel.RoomInfo,
	playSettingGroup *model2.PlaySettingGroup,
	gameInfo *model2.GameInfo,
	pickBoard *model2.PickBoard,
	board *model2.Board,
	capturedPieces *model2.CapturedPieces,
	actionHintGroup *model2.ActionHintGroup,

	msgBus core.MsgBus,
) *ResyncService {
	return &ResyncService{
		gameController:   gameController,
		userGroup:        userGroup,
		buddyGroup:       buddyGroup,
		playerGroup:      playerGroup,
		roomInfo:         roomInfo,
		playSettingGroup: playSettingGroup,
		gameInfo:         gameInfo,
		pickBoard:        pickBoard,
		board:            board,
		msgBus:           msgBus,
		capturedPieces:   capturedPieces,
		actionHintGroup:  actionHintGroup,
	}
}

func (s *ResyncService) Send(uid core.Uid) {
	gameState := s.gameController.CurrentState().ToProto(uid).(*gamegrpc.GameState)
	modelMsg := &gamegrpc.Model{
		UserGroup:       s.userGroup.ToProto(),
		RoomInfo:        s.roomInfo.ToProto(),
		BuddyGroup:      s.buddyGroup.ToProto(),
		ResyncGameState: gameState,
		GameInfo:        s.gameInfo.ToProto(),
		PlayerGroup:     s.playerGroup.ToProto(),
		// PlaySetting: see below
		Board:          s.board.ToProto(),
		CapturedPieces: s.capturedPieces.ToProto(),
		// ClaimDraw: see below
		// Surrender: see below
		PickResult:      s.pickBoard.ToProto(),
		ActionHintGroup: s.actionHintGroup.ToProto(),
	}

	if len(s.actionHintGroup.ClaimDraws) > 0 {
		claim, _ := lo.Last(s.actionHintGroup.ClaimDraws)

		modelMsg.ClaimDraw = &gamegrpc.ClaimDraw{
			ClaimUid:  claim.ClaimUid.String(),
			ClaimTurn: int32(claim.ClaimTurn),
			Claimed:   true,
			Answered:  claim.IsAccepted,
		}
	}

	surrenderUids := lo.Keys(lo.PickBy(s.actionHintGroup.Data,
		func(u core.Uid, actionHint *model2.ActionHint) bool {
			return len(actionHint.SurrenderTurns) > 0
		}))

	if len(surrenderUids) > 0 {
		modelMsg.Surrender = &gamegrpc.Surrender{
			Uid: surrenderUids[0].String(),
		}
	}

	if mySetting, ok := s.playSettingGroup.Data[uid]; ok {
		modelMsg.PlaySetting = mySetting.ToProto()
	}

	s.msgBus.Unicast(uid, core.GameStateTopic,
		gameState,
	)

	s.msgBus.Unicast(uid, core.ModelTopic,
		modelMsg,
	)
}
