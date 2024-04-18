package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx Context)

var _ Server = &HTTPServer{}

type Server interface {
	http.Handler
	Start(addr string) error
	Handle(method string, path string, handlers ...HandleFunc)
}

type HTTPServer struct {
}

//type HTTPSServer struct {
//	HTTPServer
//}

func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.server(ctx)
}

func (s *HTTPServer) server(ctx *Context) {

}

func (s *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 这里可以执行after start的操作
	return http.Serve(l, s)
}

func (s *HTTPServer) Handle(method string, path string, handlers ...HandleFunc) {
	if len(handlers) == 0 {
		panic("the number of handlers must greater than 0")
	}
}
