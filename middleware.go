package web

import (
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"time"
)

const instrumentationName = "github.com/KNICEX/go-web/web"

func Tracer() HandleFunc {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	return func(ctx *Context) {
		reqCtx := ctx.Req.Context()
		reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))
		reqCtx, span := tracer.Start(reqCtx, ctx.Req.URL.Path)
		defer span.End()

		span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
		span.SetAttributes(attribute.String("http.host", ctx.Req.Host))
		span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
		span.SetName(ctx.MatchedRoute)

	}
}

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
