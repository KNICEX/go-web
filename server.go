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

type EngineOption func(*Engine)

type Engine struct {
	*router
	*RouterGroup
}

func NewEngine(opts ...EngineOption) *Engine {
	res := &Engine{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
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
	ctx.MatchedRoute = info.node.route
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
	if len(handlers) == 0 || handlers[0] == nil {
		panic("HandleFunc is empty")
	}
	s.router.addRoute(method, path, handlers...)
}
