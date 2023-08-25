package lr

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	//// 接口声明的方式
	//var h Server = NewHTTPServer("tcp", ":8080")
	//if err := h.Server(); err != nil {
	//	panic(err)
	//}

	h := NewHTTPServer("tcp", ":8080")
	h.GET("/user/profile", func(ctx *Context) {
		fmt.Println("这是一个测试程序")
	})
	if err := h.Server(); err != nil {
		panic(err)
	}
}
