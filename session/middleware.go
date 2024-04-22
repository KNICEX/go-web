package session

import (
	"github.com/KNICEX/go-web"
	"net/http"
)

func NeedSession(m *Manager, lossSessHandler web.HandleFunc) web.HandleFunc {
	if lossSessHandler == nil {
		lossSessHandler = func(ctx *web.Context) {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
		}
	}
	return func(ctx *web.Context) {
		sess, err := m.GetSession(ctx)
		if err != nil {
			lossSessHandler(ctx)
			return
		}

		ctx.Next()
		if sess.Modified() {
			err = m.SaveSession(ctx, sess)
			if err != nil {
				ctx.Status(http.StatusInternalServerError)
				ctx.Abort()
			}
		}
	}
}
