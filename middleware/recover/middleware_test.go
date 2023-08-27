package recover

import (
	"github.com/liquanhui-99/lr"
	"log"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatusCode: http.StatusInternalServerError,
		Data:       []byte("服务器发生错误了"),
		LogFunc: func(ctx *lr.Context) {
			log.Fatalf("发生panic了，路径为: %s, status: %d, respData: %s\n",
				ctx.Req.URL.String(), ctx.Status, string(ctx.RespData))
		},
	}
	h := lr.NewHTTPServer("tcp", ":8083", lr.Use(builder.Build()))
	h.GET("/user", func(ctx *lr.Context) {
		panic("发生panic了")
	})

	if err := h.Server(); err != nil {
		panic(err)
	}
}
