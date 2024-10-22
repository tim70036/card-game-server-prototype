package model

import (
	txpokermodel "card-game-server-prototype/pkg/game/txpoker/model"
)

type PlaySettingGroup struct {
	*txpokermodel.PlaySettingGroup
}

func ProvidePlaySettingGroup() *PlaySettingGroup {
	return &PlaySettingGroup{
		PlaySettingGroup: txpokermodel.ProvidePlaySettingGroup(),
	}
}
