package server

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideGrpcServer,
	ProvideAuthInterceptor,
	ProvideLogInterceptor,
	ProvideConnectionServiceServer,
	ProvidePeerFactory,
	ProvideChatServiceServer,
	ProvideEmoteServiceServer,
)
