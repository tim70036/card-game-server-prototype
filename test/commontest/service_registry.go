package commontest

import (
	commonservice "card-game-server-prototype/pkg/common/service"
)

type serviceRegistry struct {
	baseUserService *commonservice.BaseUserService
}

func ProvideServiceRegistry(
	baseUserService *commonservice.BaseUserService,
) *serviceRegistry {
	return &serviceRegistry{
		baseUserService: baseUserService,
	}
}
