package model

import (
	txpokermodel "card-game-server-prototype/pkg/game/txpoker/model"
)

type TableProfitsGroup struct {
	*txpokermodel.TableProfitsGroup
}

func ProvideTableProfitsGroup() *TableProfitsGroup {
	return &TableProfitsGroup{
		TableProfitsGroup: txpokermodel.ProvideTableProfitsGroup(),
	}
}
