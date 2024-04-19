package web

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestRouterGroup_Group(t *testing.T) {
	base := RouterGroup{
		basePath: "/base",
		engine:   NewEngine(),
	}

	base.Use(func(ctx *Context) {
		fmt.Println("this is first middleware of base group")
	}, func(ctx *Context) {
		fmt.Println("this is second middleware of base group")
	})

	base.GET("/user", func(ctx *Context) {
		fmt.Println("GET /user handler1")
	})

	post := base.Group("/post")
	{
		post.Use(func(ctx *Context) {
			fmt.Println("this is first middleware of post group")
		})

		post.GET("/detail", func(ctx *Context) {
			fmt.Println("GET /post/detail handler1")

		})
	}
	mockRequest, err := http.NewRequest(http.MethodGet, "/base/user", nil)
	require.NoError(t, err)
	base.engine.ServeHTTP(nil, mockRequest)

	mockRequest, err = http.NewRequest(http.MethodGet, "/base/post/detail", nil)
	require.NoError(t, err)
	base.engine.ServeHTTP(nil, mockRequest)
}
