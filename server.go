package lorm

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

var _ Server = (*HTTPServer)(nil)

// Server 核心API
type Server interface {
	http.Handler
	Start(addr string) error
	addRouter(pattern, path string, handle HandleFunc)
}

type HTTPServer struct {
	*router
	// 中间件组成的链
	middlewareChain []Middleware
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

	root := h.serve
	for i := len(h.middlewareChain) - 1; i >= 0; i-- {
		root = h.middlewareChain[i](root)
	}

	root(ctx)
}

func (h *HTTPServer) serve(ctx *Context) {
	//  处理路由查找和框架的逻辑
	node, ok := h.matchRouter(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || node.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("404 Not Found"))
	}
	node.n.handler(ctx)
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
	h.addRouter(http.MethodGet, path, handle)
}

func (h *HTTPServer) POST(path string, handle HandleFunc) {
	h.addRouter(http.MethodPost, path, handle)
}

func (h *HTTPServer) PUT(path string, handle HandleFunc) {
	h.addRouter(http.MethodPut, path, handle)
}

func (h *HTTPServer) PATCH(path string, handle HandleFunc) {
	h.addRouter(http.MethodPatch, path, handle)
}

func (h *HTTPServer) DELETE(path string, handle HandleFunc) {
	h.addRouter(http.MethodDelete, path, handle)
}

func (h *HTTPServer) OPTIONS(path string, handle HandleFunc) {
	h.addRouter(http.MethodOptions, path, handle)
}
