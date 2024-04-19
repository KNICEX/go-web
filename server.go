package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

var _ Server = &Engine{}

type Server interface {
	http.Handler
	Start(addr string) error
	Handle(method string, path string, handlers ...HandleFunc)
}

type Engine struct {
	*router
}

func NewEngine() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (s *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.serve(ctx)
}

func (s *Engine) serve(ctx *Context) {
	info, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || info.node.handlers == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("404 NOT FOUND"))
		return
	}
	ctx.PathParams = info.pathParams
	for _, h := range info.node.handlers {
		h(ctx)
	}
}

func (s *Engine) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 这里可以执行after start的操作
	return http.Serve(l, s)
}

func (s *Engine) Handle(method string, path string, handlers ...HandleFunc) {
	if len(handlers) == 0 {
		panic("the number of handlers must greater than 0")
	}
	s.router.addRoute(method, path, handlers...)
}

func (s *Engine) GET(path string, handlers ...HandleFunc) {
	s.Handle(http.MethodGet, path, handlers...)
}

func (s *Engine) POST(path string, handlers ...HandleFunc) {
	s.Handle(http.MethodPost, path, handlers...)
}

func (s *Engine) PUT(path string, handlers ...HandleFunc) {
	s.Handle(http.MethodPut, path, handlers...)
}

func (s *Engine) DELETE(path string, handlers ...HandleFunc) {
	s.Handle(http.MethodDelete, path, handlers...)
}

func (s *Engine) PATCH(path string, handlers ...HandleFunc) {
	s.Handle(http.MethodPatch, path, handlers...)
}

func (s *Engine) OPTIONS(path string, handlers ...HandleFunc) {
	s.Handle(http.MethodOptions, path, handlers...)
}
