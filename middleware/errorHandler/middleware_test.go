package errorHandler

import (
	"github.com/liquanhui-99/lr"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := NewMiddlewareBuilder().
		AddCode(http.StatusNotFound).
		AddCode(http.StatusInternalServerError)
	h := lr.NewHTTPServer("tcp", ":8084", lr.Use(builder.Build()))
	h.GET("/user", func(ctx *lr.Context) {
		ctx.Status = http.StatusNotFound
	})

	h.GET("/user/profile", func(ctx *lr.Context) {
		ctx.Status = http.StatusInternalServerError
	})

	if err := h.Server(); err != nil {
		panic(err)
	}
}
