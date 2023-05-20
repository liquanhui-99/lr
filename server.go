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
}

// ServeHTTP the entry point for http requests.
func (H *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	return
}

// Start the method to start http server, received a string parameter named addr.
func (H *HTTPServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(l, nil)
}

func (H *HTTPServer) AddRouter(method string, path string, handler HandleFunc) {
	//TODO implement me
	panic("implement me")
}
