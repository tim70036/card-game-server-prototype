package service

type EventService interface {
	EvalRoundEvents() error
	EvalGameEvents() error
	Submit() error
}
