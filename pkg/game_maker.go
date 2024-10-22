package main

import (
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/game/txpoker"
	"card-game-server-prototype/pkg/util"
	"flag"
	"fmt"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GameMaker struct {
	cfg     *config.Config
	logCFG  *config.LogConfig
	testCFG *config.TestConfig
	agones  *Agones
	logger  *zap.Logger
}

func ProvideGameMaker(
	cfg *config.Config,
	logCFG *config.LogConfig,
	testCFG *config.TestConfig,
	agones *Agones,
	loggerFactory *util.LoggerFactory,
) *GameMaker {
	return &GameMaker{
		cfg:     cfg,
		logCFG:  logCFG,
		testCFG: testCFG,
		agones:  agones,
		logger:  loggerFactory.Create("GameMaker"),
	}
}

func (m *GameMaker) Run() {
	m.logger.Info("start Game",
		zap.Object("cfg", m.cfg),
		zap.Object("logCFG", m.logCFG),
		zap.Object("testCFG", m.testCFG),
	)

	if !*m.testCFG.LocalMode {
		m.logger.Info("waiting to be allocated by agones")
		go m.agones.Run()
		defer m.agones.Shutdown()
		gameServerInfo := <-m.agones.WaitAllocated
		allocationInfo := lo.Assign(gameServerInfo.ObjectMeta.Labels, gameServerInfo.ObjectMeta.Annotations)
		m.logger.Info("allocated by agones", zap.Any("allocationInfo", allocationInfo))

		allocationInfoKeys := []string{
			"room-id",
			"short-room-id",
			"game-type",
			"game-mode",
			"game-meta-uid",
			"valid-users",
		}

		for _, key := range allocationInfoKeys {
			if _, ok := allocationInfo[key]; !ok {
				m.logger.Fatal("missing field in allocation info", zap.String("missingKey", key), zap.Any("allocationInfo", allocationInfo))
			}
			if err := flag.Set(key, allocationInfo[key]); err != nil {
				m.logger.Fatal("failed to set flag using allocation info", zap.Error(err), zap.String("key", key), zap.Any("allocationInfo", allocationInfo))
			}
		}

	}

	m.logger.Info("build and run game",
		zap.Object("cfg", m.cfg),
		zap.Object("testCFG", m.testCFG),
		zap.String("LogLevel", util.LoggerLevel.Level().String()),
	)

	util.LoggerLevel.SetLevel(m.logCFG.GetLevel(string(*m.cfg.GameType)))

	switch *m.cfg.GameType {
	case gametype.TXPoker:
		var (
			txPoker *txpoker.TXPoker
			err     error
		)

		if *m.testCFG.LocalMode {
			txPoker, err = BuildLocalModeTXPoker()
		} else if *m.cfg.GameMode == gamemode.Club {
			txPoker, err = BuildClubModeTXPoker()
		} else if *m.cfg.GameMode == gamemode.Buddy {
			txPoker, err = BuildBuddyModeTXPoker()
		} else {
			txPoker, err = BuildCommonModeTXPoker()
		}

		if err != nil {
			m.logger.Fatal("failed to build TXPoker", zap.Error(err))
		}
		txPoker.Run()

	case gametype.DarkChess:
		game, err := buildDarkChess(*m.testCFG.LocalMode, *m.cfg.GameMode)
		if err != nil {
			m.logger.Fatal(fmt.Sprintf("failed to build %s", m.cfg.GameType.String()), zap.Error(err))
		}
		game.Run()

	case gametype.ZoomTXPoker:
		game, err := buildZoomTXPoker(*m.testCFG.LocalMode, *m.cfg.GameMode)
		if err != nil {
			m.logger.Fatal(fmt.Sprintf("failed to build %s", m.cfg.GameType.String()), zap.Error(err))
		}

		game.Run()

	default:
		m.logger.Fatal("unknown game type", zap.String("gameType", m.cfg.GameType.String()))
	}

	util.LoggerLevel.SetLevel(zapcore.InfoLevel)
	m.logger.Info("game end, exiting")
}
