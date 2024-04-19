package web

import (
	"net/http"
	"path"
)

type IRouterGroup interface {
	Group(relativePath string) IRouterGroup
	Use(middlewares ...HandleFunc) IRouterGroup
	Handle(httpMethod, path string, handlers ...HandleFunc) IRouterGroup
	GET(path string, handlers ...HandleFunc) IRouterGroup
	POST(path string, handlers ...HandleFunc) IRouterGroup
	DELETE(path string, handlers ...HandleFunc) IRouterGroup
	PUT(path string, handlers ...HandleFunc) IRouterGroup
	PATCH(path string, handlers ...HandleFunc) IRouterGroup
	OPTIONS(path string, handlers ...HandleFunc) IRouterGroup
}

var _ IRouterGroup = &RouterGroup{}

type RouterGroup struct {
	engine   *Engine
	handlers []HandleFunc
	basePath string
}

func (g *RouterGroup) Group(relativePath string) IRouterGroup {
	return &RouterGroup{
		engine:   g.engine,
		handlers: g.handlers,
		basePath: g.resolvePath(relativePath),
	}
}

func (g *RouterGroup) resolvePath(relativePath string) string {
	absolutePath := path.Join(g.basePath, relativePath)
	return absolutePath
}

func (g *RouterGroup) Use(middlewares ...HandleFunc) IRouterGroup {
	g.handlers = append(g.handlers, middlewares...)
	return g
}

func (g *RouterGroup) Handle(httpMethod, path string, handlers ...HandleFunc) IRouterGroup {
	if len(handlers) == 0 || handlers[0] == nil {
		panic("HandleFunc is empty")
	}
	absolutePath := g.resolvePath(path)
	combinedHandlers := append(g.handlers, handlers...)
	g.engine.addRoute(httpMethod, absolutePath, combinedHandlers...)
	return g
}

func (g *RouterGroup) GET(path string, handlers ...HandleFunc) IRouterGroup {
	return g.Handle(http.MethodGet, path, handlers...)
}

func (g *RouterGroup) POST(path string, handlers ...HandleFunc) IRouterGroup {
	return g.Handle(http.MethodPost, path, handlers...)
}

func (g *RouterGroup) DELETE(path string, handlers ...HandleFunc) IRouterGroup {
	return g.Handle(http.MethodDelete, path, handlers...)
}

func (g *RouterGroup) PUT(path string, handlers ...HandleFunc) IRouterGroup {
	return g.Handle(http.MethodPut, path, handlers...)
}

func (g *RouterGroup) PATCH(path string, handlers ...HandleFunc) IRouterGroup {
	return g.Handle(http.MethodPatch, path, handlers...)
}

func (g *RouterGroup) OPTIONS(path string, handlers ...HandleFunc) IRouterGroup {
	return g.Handle(http.MethodOptions, path, handlers...)
}
