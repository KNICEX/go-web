//go:build e2e

package session

import (
	"context"
	"github.com/KNICEX/go-web"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestNeedSession(t *testing.T) {
	e := web.NewEngine()
	e.Use(web.LoggerBuilder{}.Build(), web.RecoverBuilder{
		LogStack: true,
	}.Build())

	//sessManager := &Manager{
	//	Propagator: NewCookiePropagator(),
	//	Store:      NewMemoStore(time.Minute * 1),
	//}

	cmd := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	require.NoError(t, cmd.Ping(context.Background()).Err())
	sessManager := &Manager{
		Propagator: NewCookiePropagator(),
		Store:      NewRedisStore(cmd, time.Minute*1),
	}

	e.GET("/login/:name", func(ctx *web.Context) {
		name := ctx.Param("name")
		initSession, err := sessManager.InitSession(ctx, "session_key"+name)
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
		}

		_ = initSession.Set("name", name)
		_ = sessManager.SaveSession(ctx, initSession)

		_ = ctx.String(http.StatusOK, "login success")
	})

	e.GET("/logout", func(ctx *web.Context) {
		_, err := sessManager.GetSession(ctx)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		_ = sessManager.RemoveSession(ctx)

		_ = ctx.String(http.StatusOK, "logout success")
	})

	user := e.Group("/user")
	user.Use(NeedSession(sessManager, nil))

	{
		user.GET("/hello", func(ctx *web.Context) {
			sess, _ := sessManager.GetSession(ctx)
			name, err := sess.Get("name")
			if err != nil {
				panic(err)
			}
			_ = ctx.String(http.StatusOK, "hello "+name.(string))
		})
	}

	err := e.Start(":8080")
	t.Log(err)
}
