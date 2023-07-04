package lorm

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx Context)

var _ Server = (*HTTPServer)(nil)

// Server 核心API
type Server interface {
	http.Handler
	Start(addr string) error
	AddRouter(pattern, path string, handle HandleFunc)
}

type HTTPServer struct {
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	h.serve(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	// TODO 处理路由查找和框架的逻辑
}

func (h *HTTPServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// 这里可以注册after start
	return http.Serve(listener, h)
}

// AddRouter 注册路由，handle不提供多个，只允许注册一个，如果需要处理多个，需要用户在一个handle中实现
func (h *HTTPServer) AddRouter(pattern, path string, handle HandleFunc) {
	//TODO 注册到路由树
	panic("implement me")
}

func (h *HTTPServer) Get(path string, handle HandleFunc) {
	h.AddRouter(http.MethodGet, path, handle)
}

func (h *HTTPServer) POST(path string, handle HandleFunc) {
	h.AddRouter(http.MethodPost, path, handle)
}

func (h *HTTPServer) PUT(path string, handle HandleFunc) {
	h.AddRouter(http.MethodPut, path, handle)
}

func (h *HTTPServer) PATCH(path string, handle HandleFunc) {
	h.AddRouter(http.MethodPatch, path, handle)
}

func (h *HTTPServer) DELETE(path string, handle HandleFunc) {
	h.AddRouter(http.MethodDelete, path, handle)
}

func (h *HTTPServer) OPTIONS(path string, handle HandleFunc) {
	h.AddRouter(http.MethodOptions, path, handle)
}
