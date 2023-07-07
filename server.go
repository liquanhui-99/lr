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
	*router
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
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
