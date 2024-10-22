package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap/zapcore"
	"math"
	"time"
)

// 不是很重要的功能，不必太在意 delete expired data 時機的準確性。

type TableProfits struct {
	Uid             core.Uid
	Name            string
	CountGames      int
	SumBuyInChips   int
	SumWinLoseChips int
	expireTimestamp int64
}

func (s *TableProfits) ToProto() *txpokergrpc.TableProfits {
	return &txpokergrpc.TableProfits{
		Uid:             s.Uid.String(),
		Username:        s.Name,
		CountGames:      int32(s.CountGames),
		SumBuyInChips:   int32(s.SumBuyInChips),
		SumWinLoseChips: int32(s.SumWinLoseChips),
	}
}

func (s *TableProfits) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", s.Uid.String())
	enc.AddString("Name", s.Name)
	enc.AddInt("CountGames", s.CountGames)
	enc.AddInt("SumBuyInChips", s.SumBuyInChips)
	enc.AddInt("SumWinLoseChips", s.SumWinLoseChips)
	return nil
}

func (s *TableProfits) resetExpireTime() {
	s.expireTimestamp = time.Now().Add(constant.CacheExpirationDuration).Unix()
}

func (s *TableProfits) isExpired(now time.Time) bool {
	return s.expireTimestamp < now.Unix()
}

type TableProfitsGroup struct {
	data map[core.Uid]*TableProfits
}

func ProvideTableProfitsGroup() *TableProfitsGroup {
	return &TableProfitsGroup{
		data: make(map[core.Uid]*TableProfits),
	}
}

func (g *TableProfitsGroup) ToProto(mustHaveUids core.UidList) *txpokergrpc.TableProfitsGroup {
	msg := &txpokergrpc.TableProfitsGroup{TableProfits: make(map[string]*txpokergrpc.TableProfits)}

	for _, uid := range mustHaveUids {
		if tableProfits, ok := g.data[uid]; ok {
			msg.TableProfits[uid.String()] = tableProfits.ToProto()
		}
	}

	for uid, tableProfits := range g.data {
		if len(msg.TableProfits) >= constant.TableProfitRenderMaxSize {
			break
		}

		if _, ok := msg.TableProfits[uid.String()]; ok {
			continue
		}

		msg.TableProfits[uid.String()] = tableProfits.ToProto()
	}

	return msg
}

func (g *TableProfitsGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, stats := range g.data {
		_ = enc.AddObject(uid.String(), stats)
	}
	return nil
}

func (g *TableProfitsGroup) Get(uid core.Uid) (*TableProfits, bool) {
	d, ok := g.data[uid]
	return d, ok
}

func (g *TableProfitsGroup) Save(data *TableProfits) {
	// update
	if _, ok := g.data[data.Uid]; ok {
		g.data[data.Uid].Name = data.Name
		g.data[data.Uid].CountGames = data.CountGames
		g.data[data.Uid].SumBuyInChips = data.SumBuyInChips
		g.data[data.Uid].SumWinLoseChips = data.SumWinLoseChips
		g.data[data.Uid].resetExpireTime()
		return
	}

	// Remove the one with the oldest expiry time if the size exceeds the limit
	if len(g.data) >= constant.TableProfitsMaxSize {
		var oldestUid core.Uid
		var oldestExpireTimestamp int64 = math.MaxInt64

		for uid, d := range g.data {
			if d.expireTimestamp < oldestExpireTimestamp {
				oldestUid = uid
				oldestExpireTimestamp = d.expireTimestamp
			}
		}

		delete(g.data, oldestUid)
	}

	// Add
	g.data[data.Uid] = &TableProfits{
		Uid:             data.Uid,
		Name:            data.Name,
		CountGames:      data.CountGames,
		SumBuyInChips:   data.SumBuyInChips,
		SumWinLoseChips: data.SumWinLoseChips,
	}
	g.data[data.Uid].resetExpireTime()
}

func (g *TableProfitsGroup) EvictExpired() {
	now := time.Now()

	for uid, d := range g.data {
		if d.isExpired(now) {
			delete(g.data, uid)
		}
	}
}
