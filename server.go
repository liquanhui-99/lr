package lorm

import (
	"net"
	"net/http"
)

type HandleFunc = func(ctx Context)

type Server interface {
	// Handler the http package interface included a HTTPServer method
	http.Handler
	// Start the method to start http server
	Start(addr string) error

	// AddRouter add route to http server. method is the http method,
	// path is the routing path and handler is the method to handle request.
	AddRouter(method string, path string, handler HandleFunc)
}

var _ Server = (*HTTPServer)(nil)

// HTTPServer http server struct
type HTTPServer struct {
	router router
}

// ServeHTTP the entry point for http requests.
func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	return
}

// Start the method to start http server, received a string parameter named addr.
func (h *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	// TODO 可以做一些生命周期的管控
	return http.Serve(l, nil)
}

// AddRouter 添加路由信息
func (h *HTTPServer) AddRouter(method string, path string, handler HandleFunc) {
	h.router.AddRouter(method, path, handler)
	return
}
