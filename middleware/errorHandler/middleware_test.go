package errorHandler

import (
	"github.com/liquanhui-99/lr"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := NewMiddlewareBuilder().AddCode(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1>404 页面走丢了</h1>
	</body>
</html>
`)).AddCode(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1>401 参数错误</h1>
	</body>
</html>
`))
	h := lr.NewHTTPServer("tcp", ":8084", lr.Use(builder.Build()))
	h.GET("/user", func(ctx *lr.Context) {
		ctx.Status = http.StatusNotFound
	})
	h.GET("/user/profile", func(ctx *lr.Context) {
		ctx.Status = http.StatusBadRequest
	})
	if err := h.Server(); err != nil {
		panic(err)
	}
}
