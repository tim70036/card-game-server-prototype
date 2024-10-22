package model

import (
	txpokermodel "card-game-server-prototype/pkg/game/txpoker/model"
)

type ForceBuyInGroup struct {
	*txpokermodel.ForceBuyInGroup
}

func ProvideForceBuyInGroup() *ForceBuyInGroup {
	return &ForceBuyInGroup{
		ForceBuyInGroup: txpokermodel.ProvideForceBuyInGroup(),
	}
}
