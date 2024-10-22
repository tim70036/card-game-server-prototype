package actor

import (
	"fmt"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type AiActorFactory struct {
	loggerFactory    *util.LoggerFactory
	playSettingGroup *model2.PlaySettingGroup
	boardService     *service.BoardService
	board            *model2.Board
	actionHintGroup  *model2.ActionHintGroup
}

func ProvideAiActorFactory(
	loggerFactory *util.LoggerFactory,
	playSettingGroup *model2.PlaySettingGroup,
	boardService *service.BoardService,
	board *model2.Board,
	actionHintGroup *model2.ActionHintGroup,
) *AiActorFactory {
	return &AiActorFactory{
		loggerFactory:    loggerFactory,
		playSettingGroup: playSettingGroup,
		boardService:     boardService,
		board:            board,
		actionHintGroup:  actionHintGroup,
	}
}

func (f *AiActorFactory) Create(uid core.Uid) *AiActor {
	return &AiActor{
		logger:           f.loggerFactory.Create(fmt.Sprintf("AiActor[%s]", uid)),
		uid:              uid,
		reqDelay:         500 * time.Millisecond,
		playSettingGroup: f.playSettingGroup,
		boardService:     f.boardService,
		board:            f.board,
		actionHintGroup:  f.actionHintGroup,
	}
}

var _ Actor = (*AiActor)(nil)

type AiActor struct {
	logger           *zap.Logger
	uid              core.Uid
	reqDelay         time.Duration
	playSettingGroup *model2.PlaySettingGroup
	boardService     *service.BoardService
	board            *model2.Board
	actionHintGroup  *model2.ActionHintGroup
	playerGroup      *model2.PlayerGroup
}

func (a *AiActor) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", a.uid.String())
	enc.AddString("actorType", "AiActor")
	return nil
}

func (a *AiActor) Uid() core.Uid { return a.uid }

func (a *AiActor) DecideSkipScoreboard() (core.ActorRequestList, error) {
	return []*core.ActorRequest{
		{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &gamegrpc.SkipScoreboardRequest{},
			},
		},
	}, nil
}

func (a *AiActor) DecidePick() (core.ActorRequestList, error) {
	return []*core.ActorRequest{
		{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &gamegrpc.PickRequest{},
			},
		},
	}, nil
}

func (a *AiActor) DecideAction() (core.ActorRequestList, error) {

	// 代打順序：翻棋>移棋>吃棋>投降

	if x, y, ok := a.boardService.RandomPickUnrevealedPiece(); ok {
		return core.ActorRequestList{&core.ActorRequest{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &gamegrpc.RevealRequest{
					GridPosition: &gamegrpc.GridPosition{
						X: int32(x),
						Y: int32(y),
					},
				},
			},
		}}, nil
	}

	if x, y, toX, toY, ok := a.boardService.RandomMovePiece(a.uid); ok {
		return core.ActorRequestList{&core.ActorRequest{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &gamegrpc.MoveRequest{
					MovePiece: a.board.Cells[x][y].Piece.ToProto(),
					GoTo: &gamegrpc.GridPosition{
						X: int32(toX),
						Y: int32(toY),
					},
				},
			},
		}}, nil
	}

	if x, y, cX, cY, ok := a.boardService.RandomCapturePiece(a.uid); ok {
		return core.ActorRequestList{&core.ActorRequest{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &gamegrpc.CaptureRequest{
					MovePiece:     a.board.Cells[x][y].Piece.ToProto(),
					CapturedPiece: a.board.Cells[cX][cY].Piece.ToProto(),
					GoTo: &gamegrpc.GridPosition{
						X: int32(cX),
						Y: int32(cY),
					},
				},
			},
		}}, nil
	}

	return core.ActorRequestList{&core.ActorRequest{
		Delay: a.reqDelay,
		Req: &core.Request{
			Uid: a.uid,
			Msg: &gamegrpc.SurrenderRequest{
				Uid: a.Uid().String(),
			},
		},
	}}, nil
}
