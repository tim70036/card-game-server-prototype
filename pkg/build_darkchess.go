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
	darkchess2 "card-game-server-prototype/pkg/game/darkchess"
	api2 "card-game-server-prototype/pkg/game/darkchess/api"
	service2 "card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/util"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var providerDarkChessSet = wire.NewSet(
	config.ProviderSet,
	core.ProviderSet,
	common.ProviderSet,
	util.ProviderSet,
	wire.Value([]zap.Field{}),
	darkchess2.ProviderSet,
)

func buildDarkChess(isLocalMode bool, gameMode gamemode.GameMode) (*darkchess2.DarkChess, error) {
	if isLocalMode {
		return buildLocalModeDarkChess()
	}

	switch gameMode {
	case gamemode.Buddy:
		return buildBuddyModeDarkChess()
	case gamemode.Common:
		return buildCommonModeDarkChess()
	default:
		return nil, errors.New("invalid game mode")
	}
}

func buildLocalModeDarkChess() (*darkchess2.DarkChess, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BaseGameRepoService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.LocalUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.LocalRoomAPI)),
		wire.Bind(new(service2.EventService), new(*service2.BaseEventService)),
		wire.Bind(new(api2.GameAPI), new(*api2.LocalGameApi)),
		providerDarkChessSet,
	))
	return nil, nil
}

func buildCommonModeDarkChess() (*darkchess2.DarkChess, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BaseGameRepoService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.BaseUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),
		wire.Bind(new(service2.EventService), new(*service2.BaseEventService)),
		wire.Bind(new(api2.GameAPI), new(*api2.BaseGameApi)),
		providerDarkChessSet,
	))
	return nil, nil
}

func buildBuddyModeDarkChess() (*darkchess2.DarkChess, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BuddyModeUserService)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BuddyModeGameRepoService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.BaseUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),
		wire.Bind(new(service2.EventService), new(*service2.BaseEventService)),
		wire.Bind(new(api2.GameAPI), new(*api2.BaseGameApi)),
		providerDarkChessSet,
	))
	return nil, nil
}
