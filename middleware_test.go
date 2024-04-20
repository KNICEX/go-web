package web

import (
	"net/http"
	"testing"
)

func TestTracer(t *testing.T) {
	engine := NewEngine()
	engine.Use(Tracer())
	engine.GET("/user", func(ctx *Context) {
		_ = ctx.RespString(http.StatusOK, "user")
	})

	engine.Start(":8080")
}
