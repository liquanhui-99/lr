package lorm

import (
	"net/http"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	s := NewHTTPServer()
	// 注册路由
	s.Get("/code", func(ctx *Context) {
		ctx.Resp.WriteHeader(http.StatusOK)
		_, _ = ctx.Resp.Write([]byte("4321"))
	})

	s.Get("/query", func(ctx *Context) {
		_, err := ctx.QueryValue("test").Int()
		if err != nil {
			ctx.Resp.WriteHeader(http.StatusBadRequest)
			_, _ = ctx.Resp.Write([]byte(err.Error()))
			return
		}

		ctx.RespOKWithMessage("成功处理")
	})

	s.Get("/user/login", func(ctx *Context) {
		code := ctx.Req.URL.Query().Get("code")
		if code == "4321" {
			ctx.Resp.WriteHeader(http.StatusOK)
			_, _ = ctx.Resp.Write([]byte("登陆成功"))
		} else {
			ctx.Resp.WriteHeader(http.StatusOK)
			_, _ = ctx.Resp.Write([]byte("登陆失败"))
		}
	})

	if err := s.Start(":8080"); err != nil {
		panic(err)
	}
}
