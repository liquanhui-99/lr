package lr

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	//// 接口声明的方式
	//var h Server = NewHTTPServer("tcp", ":8080")
	//if err := h.Server(); err != nil {
	//	panic(err)
	//}

	h := NewHTTPServer("tcp", ":8080")
	h.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个before")
				fmt.Println("第三个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个before")
			}
		},
	}

	//h.GET("/user/profile", func(ctx *Context) {
	//	fmt.Println("这是一个测试程序")
	//})
	//if err := h.Server(); err != nil {
	//	panic(err)
	//}
	h.ServeHTTP(nil, &http.Request{})
}
