package commontest

type Registry struct {
	api     *apiRegistry
	service *serviceRegistry
}

func ProvideRegistry(
	api *apiRegistry,
	service *serviceRegistry,
) *Registry {
	return &Registry{
		api:     api,
		service: service,
	}
}
