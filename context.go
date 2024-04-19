package web

import "net/http"

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}

func (c *Context) Param(key string) (string, bool) {
	res, ok := c.PathParams[key]
	return res, ok
}
