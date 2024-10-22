package core

type ConnectNotification struct {
	Uid  Uid
	Done chan error
}

type EnterNotification struct {
	Uid  Uid
	Done chan error
}

type LeaveNotification struct {
	Uid  Uid
	Done chan error
}

type DisconnectNotification struct {
	Uid  Uid
	Done chan error
}

type RequestNotification struct {
	Request *Request
	Done    chan *Response
}

type TaskNotification struct {
	Task func()
}

type NextStateNotification struct {
	Trigger *StateTrigger
	Args    []any
}
