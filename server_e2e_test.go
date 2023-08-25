//go:build e2e

package lr

import (
	"net/http"
	"testing"
)

func TestServerE2e(t *testing.T) {
	//// 接口声明的方式
	//var h Server = NewHTTPServer("tcp", ":8080")
	//if err := h.Server(); err != nil {
	//	panic(err)
	//}

	h := NewHTTPServer("tcp", ":8081")

	h.POST("/user/profile", func(ctx *Context) {
		t.Log("成功")
		ctx.Resp.WriteHeader(http.StatusOK)
	})

	h.GET("/user/login/:id", func(ctx *Context) {
		id := ctx.pathParams["id"]
		t.Log("参数为：", id)
	})

	if err := h.Server(); err != nil {
		panic(err)
	}
}
