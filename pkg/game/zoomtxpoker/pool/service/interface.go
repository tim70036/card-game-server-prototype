package service

type GameRepoService interface {
	FetchGameInfo() error
	FetchGameWater() error
}
