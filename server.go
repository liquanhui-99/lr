package lr

import (
	"log"
	"net"
	"net/http"
)

type Server interface {
	// Handler 组合http的接口
	http.Handler
	// Server 启动服务的方法
	Server() error
	// AddRoute 注册路由信息
	AddRoute(method, path string, handler HandleFunc)
}

var _ Server = (*HTTPServer)(nil)

// HTTPServer http的实现
type HTTPServer struct {
	// 监听的地址
	addr string
	// 网路
	network string
	// 组合路由
	*router
	// server层面上的Middleware
	mdls []Middleware
	// 模版渲染引擎
	tplEngine TemplateEngine
}

type HTTPServerOptions func(server *HTTPServer)

func (h *HTTPServer) AddRoute(method, path string, handler HandleFunc) {
	//TODO implement me
	panic("implement me")
}

func NewHTTPServer(network, addr string, opts ...HTTPServerOptions) *HTTPServer {
	s := &HTTPServer{
		addr:    addr,
		network: network,
		router:  newRouter(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Use 传入中间件
func Use(mdls ...Middleware) HTTPServerOptions {
	return func(s *HTTPServer) {
		s.mdls = mdls
	}
}

// Template 初始化渲染模版的引擎
func Template(engine TemplateEngine) HTTPServerOptions {
	return func(s *HTTPServer) {
		s.tplEngine = engine
	}
}

// GET 注册GET方法
func (h *HTTPServer) GET(path string, handler HandleFunc) {
	h.router.addRouter(http.MethodGet, path, handler)
}

// POST 注册POST方法
func (h *HTTPServer) POST(path string, handler HandleFunc) {
	h.router.addRouter(http.MethodPost, path, handler)
}

// PUT 注册PUT方法
func (h *HTTPServer) PUT(path string, handler HandleFunc) {
	h.router.addRouter(http.MethodPut, path, handler)
}

// PATCH 注册PATCH方法
func (h *HTTPServer) PATCH(path string, handler HandleFunc) {
	h.router.addRouter(http.MethodPatch, path, handler)
}

// DELETE 注册DELETE方法
func (h *HTTPServer) DELETE(path string, handler HandleFunc) {
	h.router.addRouter(http.MethodDelete, path, handler)
}

// OPTIONS 注册DELETE方法
func (h *HTTPServer) OPTIONS(path string, handler HandleFunc) {
	h.router.addRouter(http.MethodOptions, path, handler)
}

// ServerHTTP 处理请求的入口方法
func (h *HTTPServer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:       request,
		Resp:      response,
		TplEngine: h.tplEngine,
	}

	// 中间件的处理逻辑，从后往前的方式挂载
	root := h.serve
	for i := len(h.mdls) - 1; i >= 0; i-- {
		root = h.mdls[i](root)
	}

	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			// 在这里设置响应状态码和响应的数据信息
			h.flushResp(ctx)
		}
	}
	root = m(root)

	root(ctx)
}

func (h *HTTPServer) flushResp(ctx *Context) {
	if ctx.Status != 0 {
		ctx.Resp.WriteHeader(ctx.Status)
	}

	if len(ctx.RespData) != 0 {
		n, err := ctx.Resp.Write(ctx.RespData)
		if err != nil || n != len(ctx.RespData) {
			log.Fatalf("写入响应失败 %v", err)
		}
	}
}

// serve 需要先查询路由树，执行命中的逻辑
func (h *HTTPServer) serve(ctx *Context) {
	res, ok := h.router.findRouter(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || res.n.handler == nil {
		// 不存在路径或者路径查到但是没有handler
		ctx.Status = http.StatusNotFound
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("NOT FOUND"))
		return
	}

	ctx.matchedPath = res.n.fullPath
	ctx.pathParams = res.pathParams
	res.n.handler(ctx)
}

// Server 启动程序
func (h *HTTPServer) Server() error {
	listener, err := net.Listen(h.network, h.addr)
	if err != nil {
		return err
	}

	return http.Serve(listener, h)
}

// HTTPSServer https的实现
type HTTPSServer struct {
	HTTPServer
}
