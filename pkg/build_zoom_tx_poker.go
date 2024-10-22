//go:build wireinject
// +build wireinject

package main

import (
	"errors"
	"card-game-server-prototype/pkg/common"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	pool2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool"
	api2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/api"
	service2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/service"
	"card-game-server-prototype/pkg/util"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var zoomTXPokerProviderSet = wire.NewSet(
	config.ProviderSet,
	core.ProviderSet,
	common.ProviderSet,
	util.ProviderSet,
	pool2.ProviderSet,
)

func buildZoomTXPoker(isLocalMode bool, gameMode gamemode.GameMode) (*pool2.ZoomTXPokerPool, error) {
	zapField := []zap.Field{
		zap.String("zoom-module", "zoom-pool"),
	}

	if isLocalMode {
		return buildLocalModeZoomTXPoker(zapField)
	}

	switch gameMode {
	case gamemode.Common:
		return buildCommonModeZoomTXPoker(zapField)
	default:
		return nil, errors.New("invalid game mode")
	}
}

func buildLocalModeZoomTXPoker(
	_ []zap.Field,
) (*pool2.ZoomTXPokerPool, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.LocalUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.LocalRoomAPI)),
		wire.Bind(new(api2.GameAPI), new(*api2.LocalGameAPI)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BaseGameRepoService)),
		zoomTXPokerProviderSet,
	))
	return nil, nil
}

func buildCommonModeZoomTXPoker(
	_ []zap.Field,
) (*pool2.ZoomTXPokerPool, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.BaseUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),
		wire.Bind(new(api2.GameAPI), new(*api2.BaseGameAPI)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BaseGameRepoService)),
		zoomTXPokerProviderSet,
	))
	return nil, nil
}
