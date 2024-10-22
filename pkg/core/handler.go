package core

type Handler interface {
	HandleConnect(uid Uid) error
	HandleEnter(uid Uid) error
	HandleLeave(uid Uid) error
	HandleDisconnect(uid Uid) error
	HandleRequest(request *Request) *Response
	mustEmbedBaseHandler()
}

type BaseHandler struct{}

func (handler *BaseHandler) HandleConnect(uid Uid) error {
	return nil
}

func (handler *BaseHandler) HandleEnter(uid Uid) error {
	return nil
}

func (handler *BaseHandler) HandleLeave(uid Uid) error {
	return nil
}

func (handler *BaseHandler) HandleDisconnect(uid Uid) error {
	return nil
}

func (handler *BaseHandler) HandleRequest(request *Request) *Response {
	return nil
}

func (handler *BaseHandler) mustEmbedBaseHandler() {}
