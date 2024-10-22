package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	wire.Value(ClientCFG),
	wire.Value(ServerCFG),
	wire.Value(LogCFG),
	wire.Value(CFG),
	wire.Value(APICFG),
	wire.Value(TestCFG),
	wire.Value(RedisClientCFG),
)
