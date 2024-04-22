package web

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MiddlewareBuilder interface {
	Build() HandleFunc
}

func DefaultLogFunc(info string) {
	log.Println(info)
}

type LoggerBuilder struct {
	LogFunc func(log string)
}

func (l LoggerBuilder) Build() HandleFunc {
	if l.LogFunc == nil {
		l.LogFunc = DefaultLogFunc
	}
	return func(ctx *Context) {
		startTime := time.Now()
		ctx.Next()
		defer func() {
			al := accessLog{
				Host:    ctx.Req.Host,
				Route:   ctx.MatchedRoute,
				Method:  ctx.Req.Method,
				Path:    ctx.Req.URL.Path,
				Latency: time.Since(startTime),
			}
			data, _ := json.Marshal(al)
			l.LogFunc(string(data))
		}()
	}
}

type accessLog struct {
	Host    string
	Route   string
	Method  string
	Path    string
	Latency time.Duration
}

type PrometheusBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

func (p PrometheusBuilder) Build() HandleFunc {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: p.Namespace,
		Subsystem: p.Subsystem,
		Name:      p.Name,
		Help:      p.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(vector)

	return func(ctx *Context) {
		startTime := time.Now()
		defer func() {
			pattern := ctx.MatchedRoute
			if pattern == "" {
				pattern = "unknown"
			}
			vector.WithLabelValues(pattern, ctx.Req.Method, strconv.Itoa(ctx.StatusCode)).
				Observe(float64(time.Since(startTime).Milliseconds()))
		}()
		ctx.Next()
	}
}

type RecoverBuilder struct {
	LogFunc  func(log string)
	LogStack bool
	Handler  HandleFunc
}

func DefaultRecoverHandler(ctx *Context) {
	ctx.StatusCode = 500
	ctx.RespData = []byte("Internal Server Error")
}

func (r RecoverBuilder) Build() HandleFunc {
	if r.LogFunc == nil {
		r.LogFunc = DefaultLogFunc
	}
	if r.Handler == nil {
		r.Handler = DefaultRecoverHandler
	}
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				if r.LogStack {
					r.LogFunc(trace(fmt.Sprintf("%s", err)))
				} else {
					r.LogFunc(fmt.Sprintf("%s", err))
				}
				r.Handler(ctx)
			}
		}()
		ctx.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
