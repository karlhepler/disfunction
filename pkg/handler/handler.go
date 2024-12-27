package handler

import "context"

type Interface[RequestData, ResponseData WithContext] interface {
	Handle(req RequestData, res Sender[ResponseData])
}

type Sender[ResponseData WithContext] interface {
	Send(ResponseData)
}

type WithContext interface {
	Context() context.Context
}

// Handle is a simple wrapper function that ensures hdl adheres to the
// interface. I could imagine expanding it to do more in the future, but
// this is enough for now.
//
// Since this gets run with every request, this is a great place to tap into.
func Handle[RequestData, ResponseData WithContext](hdl Interface[RequestData, ResponseData], req RequestData, res Sender[ResponseData]) {
	hdl.Handle(req, res)
}
