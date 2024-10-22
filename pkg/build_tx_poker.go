//go:build wireinject
// +build wireinject

package main

import (
	"card-game-server-prototype/pkg/common"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	txpoker2 "card-game-server-prototype/pkg/game/txpoker"
	api2 "card-game-server-prototype/pkg/game/txpoker/api"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/util"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var txPokerProviderSet = wire.NewSet(
	config.ProviderSet,
	core.ProviderSet,
	common.ProviderSet,
	util.ProviderSet,
	wire.Value([]zap.Field{}),
	txpoker2.ProviderSet,
)

func BuildLocalModeTXPoker() (*txpoker2.TXPoker, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.LocalUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.LocalRoomAPI)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BaseGameRepoService)),
		wire.Bind(new(service2.EventService), new(*service2.BaseEventService)),
		wire.Bind(new(service2.SeatStatusService), new(*service2.BaseSeatStatusService)),
		wire.Bind(new(service2.JackpotService), new(*service2.BaseJackpotService)),
		wire.Bind(new(api2.GameAPI), new(*api2.LocalGameAPI)),
		txPokerProviderSet,
	))
	return nil, nil
}

func BuildCommonModeTXPoker() (*txpoker2.TXPoker, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.BaseUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BaseGameRepoService)),
		wire.Bind(new(service2.EventService), new(*service2.BaseEventService)),
		wire.Bind(new(service2.SeatStatusService), new(*service2.BaseSeatStatusService)),
		wire.Bind(new(service2.JackpotService), new(*service2.BaseJackpotService)),
		wire.Bind(new(api2.GameAPI), new(*api2.BaseGameAPI)),
		txPokerProviderSet,
	))
	return nil, nil
}

func BuildClubModeTXPoker() (*txpoker2.TXPoker, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.ClubModeUserService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(service2.GameRepoService), new(*service2.ClubModeGameRepoService)),
		wire.Bind(new(service2.EventService), new(*service2.BaseEventService)),
		wire.Bind(new(service2.SeatStatusService), new(*service2.ClubModeSeatStatusService)),
		wire.Bind(new(service2.JackpotService), new(*service2.ClubModeJackpotService)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),
		wire.Bind(new(api2.GameAPI), new(*api2.BaseGameAPI)),
		wire.Bind(new(api2.ClubGameAPI), new(*api2.ClubModeGameAPI)),
		wire.Bind(new(commonapi.UserAPI), new(*api2.ClubModeUserAPI)),
		wire.Bind(new(api2.ClubMemberAPI), new(*api2.ClubModeMemberApi)),
		txPokerProviderSet,
	))
	return nil, nil
}

func BuildBuddyModeTXPoker() (*txpoker2.TXPoker, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonservice.UserService), new(*service2.BaseUserService)),
		wire.Bind(new(commonservice.RoomService), new(*commonservice.BaseRoomService)),
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.BaseUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),
		wire.Bind(new(service2.GameRepoService), new(*service2.BuddyModeGameRepoService)),
		wire.Bind(new(service2.EventService), new(*service2.BaseEventService)),
		wire.Bind(new(service2.SeatStatusService), new(*service2.BaseSeatStatusService)),
		wire.Bind(new(service2.JackpotService), new(*service2.BaseJackpotService)),
		wire.Bind(new(api2.GameAPI), new(*api2.BaseGameAPI)),
		txPokerProviderSet,
	))
	return nil, nil
}
