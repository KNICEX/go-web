package web

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
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
