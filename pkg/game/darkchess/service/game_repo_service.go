package service

type GameRepoService interface {
	CreateGame() error
	CreateRound() error
	SubmitRoundScore() error
	SubmitGameScore() error

	FetchRoundResult() error
	FetchGameResult() error
	FetchGameSetting() error
}
