package lorm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRouter(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		path    string
	}{
		{
			pattern: http.MethodGet,
			path:    "/",
		},
		{
			pattern: http.MethodGet,
			path:    "/login",
		},
		{
			pattern: http.MethodGet,
			path:    "/login/verifyCode",
		},
		{
			pattern: http.MethodPost,
			path:    "/api/user/login",
		},
		{
			pattern: http.MethodDelete,
			path:    "/user/code",
		},
	}
	var mockHandler HandleFunc = func(c *Context) {}
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"api": &node{
						path: "api",
						children: map[string]*node{
							"user": &node{
								path: "user",
								children: map[string]*node{
									"login": &node{
										path:     "login",
										children: map[string]*node{},
										handler:  mockHandler,
									},
								},
							},
						},
					},
				},
			},
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"login": &node{
						path: "login",
						children: map[string]*node{
							"verifyCode": &node{
								path:     "verifyCode",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
						handler: mockHandler,
					},
				},
				handler: mockHandler,
			},
			http.MethodDelete: &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"code": &node{
								path:     "code",
								children: map[string]*node{},
								handler:  mockHandler,
							},
						},
					},
				},
			},
		},
	}
	router := newRouter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router.addRouter(tc.pattern, tc.path, mockHandler)
		})
	}
	msg, ok := router.equal(wantRouter)
	if !ok {
		fmt.Println("结果为: ", msg)
	}
	fmt.Println("匹配成功！")
}

func TestPanic(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	router := newRouter()
	assert.Panicsf(t, func() {
		router.addRouter(http.MethodPut, "", mockHandler)
	}, "请求路径不能为空")
	assert.Panicsf(t, func() {
		router.addRouter(http.MethodPut, "user", mockHandler)
	}, "请求路径不是以/开头")
	assert.Panicsf(t, func() {
		router.addRouter(http.MethodPut, "/user/", mockHandler)
	}, "请求路径不能以/结尾")
	assert.Panicsf(t, func() {
		router.addRouter(http.MethodPut, "/user//code", mockHandler)
	}, "请求路径不能包含连续的/")
}

func TestRepeatedPath(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	router := newRouter()
	router.addRouter(http.MethodGet, "/", mockHandler)
	assert.Panicsf(t, func() {
		router.addRouter(http.MethodGet, "/", mockHandler)
	}, "请求路径[/]已注册")
	router.addRouter(http.MethodPut, "/user/login", mockHandler)
	assert.Panicsf(t, func() {
		router.addRouter(http.MethodPut, "/user/login", mockHandler)
	}, "请求路径[/user/login]已注册")
}

func (r *router) equal(d *router) (string, bool) {
	if r.trees == nil || d.trees == nil {
		return fmt.Sprintf("匹配失败"), false
	}
	if len(r.trees) != len(d.trees) {
		return fmt.Sprintf("匹配失败，路由树数量不匹配"), false
	}
	for method, tree := range r.trees {
		val, ok := d.trees[method]
		if !ok {
			return fmt.Sprintf("请求的方法不匹配"), false
		}
		msg, ok := tree.equal(val)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(d *node) (string, bool) {
	if n.path != d.path {
		return fmt.Sprint("节点的路径不匹配"), false
	}

	// 处理路径参数的查找
	if n.paramChild != nil {
		msg, ok := n.paramChild.equal(d.paramChild)
		return msg, ok
	}

	// 比对handler是否一致
	nHandler := reflect.ValueOf(n.handler)
	dHandler := reflect.ValueOf(d.handler)
	if nHandler != dHandler {
		return fmt.Sprintf("节点的handler不匹配"), false
	}

	if len(n.children) != len(d.children) {
		return fmt.Sprintf("子节点数量不想等"), false
	}
	for path, node := range n.children {
		val, ok := d.children[path]
		if !ok {
			return fmt.Sprintf("子节点路径不匹配"), false
		}
		msg, ok := node.equal(val)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func TestMatchRouter(t *testing.T) {
	testRouter := []struct {
		name    string
		pattern string
		path    string
	}{
		{
			name:    "",
			pattern: http.MethodGet,
			path:    "/user/login",
		},
	}

	var mockHandler HandleFunc = func(ctx *Context) {}
	router := newRouter()
	for _, tc := range testRouter {
		t.Run(tc.name, func(t *testing.T) {
			router.addRouter(tc.pattern, tc.path, mockHandler)
		})
	}

	testCases := []struct {
		name      string
		pattern   string
		path      string
		wantFound bool
		wantNode  *pathInfo
	}{
		{
			name:    "not exist method",
			pattern: http.MethodOptions,
			path:    "/user/login",
		},
		{
			name:      "match router",
			pattern:   http.MethodGet,
			path:      "/user/login",
			wantFound: true,
			wantNode: &pathInfo{
				n: &node{
					path:     "login",
					handler:  mockHandler,
					children: map[string]*node{},
				},
			},
		},
		{
			name:    "match not found",
			pattern: http.MethodPost,
			path:    "/order",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, ok := router.matchRouter(tc.pattern, tc.path)
			assert.Equal(t, tc.wantFound, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.wantNode.n.path, node.n.path)
			assert.Equal(t, tc.wantNode.n.children, node.n.children)
			yHandler := reflect.ValueOf(tc.wantNode.n.handler)
			nHandler := reflect.ValueOf(node.n.handler)
			assert.True(t, yHandler == nHandler)
		})
	}
}

func TestRouter_pathParam(t *testing.T) {
	testCases := []struct {
		pattern string
		name    string
		path    string
	}{
		{
			name:    "path parameter",
			pattern: http.MethodGet,
			path:    "/task/result/:id",
		},
		{
			name:    "path parameter",
			pattern: http.MethodGet,
			path:    "/user/:id",
		},
		{
			name:    "同路径同时存在静态路径和路径参数",
			pattern: http.MethodGet,
			path:    "/task/:id",
		},
	}
	var mockHandler HandleFunc = func(ctx *Context) {}
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path: "/",
				children: map[string]*node{
					"task": &node{
						path: "task",
						children: map[string]*node{
							"result": &node{
								path: "result",
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
						paramChild: &node{
							path:    ":id",
							handler: mockHandler,
						},
					},
					"user": &node{
						path: "user",
						paramChild: &node{
							path:    ":id",
							handler: mockHandler,
						},
					},
				},
			},
		},
	}
	r := newRouter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r.addRouter(tc.pattern, tc.path, mockHandler)
		})
	}
	msg, ok := r.equal(wantRouter)
	if !ok {
		panic(msg)
	}
}

func TestRouter_MatchParam(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	testCases := []struct {
		name     string
		pattern  string
		path     string
		wantNode *pathInfo
	}{
		{
			name:    "多个路径",
			pattern: http.MethodGet,
			path:    "/task/result/:id",
			wantNode: &pathInfo{
				n: &node{
					path:    ":id",
					handler: mockHandler,
				},
			},
		},
		{
			name:    "两个路径",
			pattern: http.MethodGet,
			path:    "/order/:name",
			wantNode: &pathInfo{
				n: &node{
					path:    ":name",
					handler: mockHandler,
				},
			},
		},
	}
	r := newRouter()
	for _, tc := range testCases {
		r.addRouter(tc.pattern, tc.path, mockHandler)
	}

	for _, tc := range testCases {
		node, ok := r.matchRouter(tc.pattern, tc.path)
		if !ok {
			fmt.Println("不匹配")
			return
		}
		assert.Equal(t, tc.wantNode.n.path, node.n.path)
		yHandler := reflect.ValueOf(tc.wantNode.n.handler)
		nHandler := reflect.ValueOf(node.n.handler)
		assert.True(t, yHandler == nHandler)
	}
}

func TestRouter_PathParam(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	r.addRouter(http.MethodGet, "/task/:taskName", mockHandler)
	testCases := []struct {
		name      string
		pattern   string
		path      string
		wantFound bool
		info      *pathInfo
	}{
		{
			name:      "获取路径参数",
			pattern:   http.MethodGet,
			path:      "/task/:testName",
			wantFound: true,
			info: &pathInfo{
				n: &node{
					path:    ":taskName",
					handler: mockHandler,
				},
				pathParams: map[string]string{
					"taskName": "testName",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, ok := r.matchRouter(tc.pattern, tc.path)
			assert.Equal(t, tc.wantFound, ok)
			assert.Equal(t, tc.info.n.path, info.n.path)
			assert.Equal(t, tc.info.pathParams, info.pathParams)
			msg, ok := tc.info.n.equal(info.n)
			if !ok {
				panic(msg)
			}
		})
	}
}

func Benchmark_MatchRouter(b *testing.B) {
	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}
	testCases := []struct {
		name    string
		pattern string
		path    string
	}{
		{
			name:    "Get method",
			pattern: http.MethodGet,
			path:    "/api/getCode",
		},
		{
			name:    "Get method multiple path",
			pattern: http.MethodGet,
			path:    "/task/result/detail",
		},
		{
			name:    "Get method multiple path1",
			pattern: http.MethodGet,
			path:    "/task/result/detail/instance",
		},
		{
			name:    "Get method multiple path",
			pattern: http.MethodGet,
			path:    "/task/result/:id",
		},
		{
			name:    "Post method",
			pattern: http.MethodPost,
			path:    "/login",
		},
		{
			name:    "Put method",
			pattern: http.MethodPut,
			path:    "/task/login/:username",
		},
		{
			name:    "Delete method",
			pattern: http.MethodPut,
			path:    "/task/login/:id",
		},
	}
	for _, tc := range testCases {
		r.addRouter(tc.pattern, tc.path, mockHandler)
	}
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			b.Run(tc.name, func(b *testing.B) {
				_, ok := r.matchRouter(tc.pattern, tc.path)
				if !ok {
					panic("未匹配到路径")
				}
			})
		}
	}
}
