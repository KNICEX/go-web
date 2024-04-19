package web

import (
	"encoding/json"
	"log"
	"time"
)

var logFunc func(log string) = func(info string) {
	log.Println(info)
}

func Logger() HandleFunc {
	return func(ctx *Context) {
		// before
		t := time.Now()
		ctx.Next()
		defer func() {
			l := accessLog{
				Host:    ctx.Req.Host,
				Route:   ctx.MatchedRoute,
				Method:  ctx.Req.Method,
				Path:    ctx.Req.URL.Path,
				Latency: time.Since(t),
			}
			data, _ := json.Marshal(l)
			logFunc(string(data))
		}()
		// after
	}
}

type accessLog struct {
	Host    string
	Route   string
	Method  string
	Path    string
	Latency time.Duration
}
