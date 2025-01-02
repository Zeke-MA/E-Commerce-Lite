package server

type HttpResponseMessages string

const (
	MsgOK            HttpResponseMessages = "OK"
	MsgInternalError HttpResponseMessages = "Internal Server Error"
	MsgConflict      HttpResponseMessages = "The resource already exists"
)
