package lr

import "net/http"

type Context struct {
	// Req 接受的请求信息
	Req *http.Request
	// Resp 返回响应
	Resp http.ResponseWriter
}
