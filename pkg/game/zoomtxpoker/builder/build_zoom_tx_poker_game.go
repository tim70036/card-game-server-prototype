//go:build wireinject
// +build wireinject

package builder

import (
	"card-game-server-prototype/pkg/common"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker"
	txpokerapi "card-game-server-prototype/pkg/game/txpoker/api"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	game2 "card-game-server-prototype/pkg/game/zoomtxpoker/game"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/api"
	service3 "card-game-server-prototype/pkg/game/zoomtxpoker/game/service"
	session2 "card-game-server-prototype/pkg/game/zoomtxpoker/game/session"
	"card-game-server-prototype/pkg/util"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var zoomTXPokerGameProviderSet = wire.NewSet(
	config.ProviderSet,
	core.ProviderSet,
	util.ProviderSet,
	common.ProviderSet,

	txpoker.ProviderSet,
	game2.ProviderSet,
)

func BuildZoomTXPokerGame(
	_ *session2.UserSessionGroup,
	_ *session2.GameInfoSession,
	_ *session2.StatsSessionGroup,
	_ *session2.ForceBuyInSessionGroup,
	_ *session2.TableProfitsSessionGroup,
	_ []zap.Field,
) (*game2.ZoomTXPokerGame, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.BaseUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BaseGameRepoService)),
		wire.Bind(new(service2.EventService), new(*service3.BaseEventService)),
		wire.Bind(new(service2.SeatStatusService), new(*service3.BaseSeatStatusService)),
		wire.Bind(new(txpokerapi.GameAPI), new(*api.BaseGameAPI)),
		wire.Bind(new(service2.JackpotService), new(*service3.BaseJackpotService)),
		zoomTXPokerGameProviderSet,
	))

	return nil, nil
}
