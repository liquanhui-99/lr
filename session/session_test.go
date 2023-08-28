package session

import (
	"github.com/liquanhui-99/lr"
	"net/http"
	"testing"
)

func TestSession(t *testing.T) {
	h := lr.NewHTTPServer("tcp", "8081", lr.Use(func(next lr.HandleFunc) lr.HandleFunc {
		return func(ctx *lr.Context) {
			if ctx.Req.URL.Path == "/login" {
				next(ctx)
			}

			manager := Manager{}
			_, err := manager.GetSession(ctx)
			if err != nil {
				ctx.Status = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登陆")
				return
			}

			// 刷新session信息
			if err = manager.Refresh(ctx); err != nil {
				t.Log(err)
			}

			next(ctx)
		}
	}))

	h.POST("/login", func(ctx *lr.Context) {
		manager := Manager{}
		_, err := manager.InitSession(ctx)
		if err != nil {
			ctx.Status = http.StatusUnauthorized
			ctx.RespData = []byte("请重新登陆")
			return
		}
		ctx.Status = http.StatusOK
		ctx.RespData = []byte("登陆成功")
	})

	h.POST("/logout", func(ctx *lr.Context) {
		manager := Manager{}
		err := manager.RemoveSession(ctx)
		if err != nil {
			ctx.Status = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.Status = http.StatusOK
		ctx.RespData = []byte("退出成功")
	})

	if err := h.Server(); err != nil {
		panic(err)
	}
}
