package accesslog

import (
	"github.com/liquanhui-99/lr"
	"log"
	"net/http"
	"testing"
)

func TestAccessLog(t *testing.T) {
	accessLog := NewMiddleBuilder().LogFunc(func(s string) {
		log.Fatalf("结果为: %s\n", s)
	})
	s := lr.NewHTTPServer("tcp", ":8081",
		lr.Use(accessLog.Build()))

	s.GET("/", func(ctx *lr.Context) {
		t.Log("成功")
		ctx.Resp.WriteHeader(http.StatusOK)
	})

	if err := s.Server(); err != nil {
		panic(err)
	}
}
