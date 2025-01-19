package server

type HttpResponseMessages string

const (
	MsgOK            HttpResponseMessages = "OK"
	MsgInternalError HttpResponseMessages = "Internal Server Error"
	MsgConflict      HttpResponseMessages = "The resource already exists"
	MsgNotFound      HttpResponseMessages = "Not Found"
	MsgUnauthorized  HttpResponseMessages = "Unauthorized Request"
	MsgBadRequest    HttpResponseMessages = "Bad Request"
)
