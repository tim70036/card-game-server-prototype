package actor

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideActorGroup,
	ProvideBaseActorFactory,
	ProvideDummyActorFactory,
)
