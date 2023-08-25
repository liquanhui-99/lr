package lr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRouter(t *testing.T) {
	testCases := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "post",
			method: http.MethodPost,
			path:   "/user/signIn",
		},
	}

	var mockHanlerFunc HandleFunc = func(ctx *Context) {}
	wantTrees := &router{
		trees: map[string]*node{
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"signIn": {
								path:     "signIn",
								children: map[string]*node{},
								handler:  mockHanlerFunc,
							},
						},
					},
				},
			},
		},
	}

	r := newRouter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r.addRouter(tc.method, tc.path, mockHanlerFunc)
		})
	}

	// 断言路由树是相等的，HandleFunc不能直接比较
	msg, ok := r.equal(*wantTrees)
	if !ok {
		t.Log(msg)
		return
	}
	t.Log("匹配成功")
}

// equal 比对路由树是否相等
func (r *router) equal(dest router) (string, bool) {
	if len(r.trees) != len(dest.trees) {
		return fmt.Sprintf("路由树不匹配"), false
	}

	for mtd, tr := range dest.trees {
		tree, ok := r.trees[mtd]
		if !ok {
			return fmt.Sprintf("路由方法不匹配"), false
		}
		// 比较node是否相等
		msg, ok := tree.equal(tr)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

// equal 比较node节点是否一致
func (n *node) equal(dst *node) (string, bool) {
	if n.path != dst.path {
		return fmt.Sprintf("路径不匹配"), false
	}

	if len(n.children) != len(dst.children) {
		return fmt.Sprintf("子节点数量不匹配"), false
	}

	// 比较handler
	nHandler := reflect.ValueOf(n.handler)
	dstHandler := reflect.ValueOf(dst.handler)
	if nHandler != dstHandler {
		return fmt.Sprintf("节点Handler不匹配"), false
	}

	// 比较子节点
	for path, child := range dst.children {
		nd, ok := n.children[path]
		if !ok {
			return fmt.Sprintf("子节点路径不匹配"), false
		}

		msg, ok := nd.equal(child)
		if !ok {
			return msg, false
		}
	}

	return "", true
}

// TestPanic_EmptyPath 测试空路径
func TestPanic_EmptyPath(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()

	assert.Panics(t, func() {
		r.addRouter(http.MethodPost, "", mockHandler)
	})
}

// TestPanic_Prefix 测试不以/开头
func TestPanic_Prefix(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()

	assert.Panics(t, func() {
		r.addRouter(http.MethodPost, "user/profile", mockHandler)
	}, "请求路径必须以/开头")
}

// TestPanic_Suffix 测试不以/结尾
func TestPanic_Suffix(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()

	assert.Panics(t, func() {
		r.addRouter(http.MethodPost, "/user/profile/", mockHandler)
	}, "请求路径不能以/结尾")

	assert.Panics(t, func() {
		r.addRouter(http.MethodGet, "/user/", mockHandler)
	}, "请求路径不能以/结尾")
}

// TestPanic_DoubleSlash 测试不能包含连续的斜杠
func TestPanic_DoubleSlash(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()

	assert.Panics(t, func() {
		r.addRouter(http.MethodGet, "/user//profile", mockHandler)
	}, "请求路径不能包含连续的/")

	assert.Panics(t, func() {
		r.addRouter(http.MethodGet, "//user/profile", mockHandler)
	}, "请求路径不能包含连续的/")

	assert.Panics(t, func() {
		r.addRouter(http.MethodGet, "/user/profile//", mockHandler)
	}, "请求路径不能以/结尾")
}

// TestPanic_DuplicateRootPath 根节点重复注册
func TestPanic_DuplicateRootPath(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	r.addRouter(http.MethodPost, "/", mockHandler)
	assert.Panics(t, func() {
		r.addRouter(http.MethodPost, "/", mockHandler)
	}, "路由冲突，重复注册[/]")
}

// TestPanic_DuplicatePath 普通路径的重复注册
func TestPanic_DuplicatePath(t *testing.T) {
	var mockHandler HandleFunc = func(ctx *Context) {}
	r := newRouter()
	r.addRouter(http.MethodPost, "/user/signIn", mockHandler)
	assert.Panics(t, func() {
		r.addRouter(http.MethodPost, "/user/signIn", mockHandler)
	}, "路由冲突，重复注册[/user/signIn]")
}

func TestFindPath(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		method string
	}{
		{
			name:   "GET root",
			path:   "/",
			method: http.MethodGet,
		},
		{
			name:   "GET",
			path:   "/user/profile",
			method: http.MethodGet,
		},
		{
			name:   "POST",
			path:   "/user/signIn",
			method: http.MethodPost,
		},
		{
			name:   "POST",
			path:   "/user/signUp",
			method: http.MethodPost,
		},
		{
			name:   "GET",
			path:   "/user/signUp/:id",
			method: http.MethodGet,
		},
		{
			name:   "GET",
			path:   "/user/signUp",
			method: http.MethodGet,
		},
	}

	r := newRouter()
	var mockHandler HandleFunc = func(ctx *Context) {}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r.addRouter(tc.method, tc.path, mockHandler)
		})
	}

	testCases1 := []struct {
		name      string
		path      string
		method    string
		wantFound bool
		wantNode  *node
	}{
		{
			name:      "GET root",
			path:      "/",
			method:    http.MethodGet,
			wantFound: true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"profile": {
								path:     "profile",
								children: map[string]*node{},
								handler:  mockHandler,
							},
							"signUp": {
								path:     "signUp",
								children: map[string]*node{},
								handler:  mockHandler,
								paramChild: &node{
									path:    ":id",
									handler: mockHandler,
								},
							},
						},
					},
				},
			},
		},
		{
			name:      "PUT",
			path:      "/user/editProfile",
			method:    http.MethodPut,
			wantFound: false,
		},
		{
			name:      "POST",
			path:      "/user/signUp",
			method:    http.MethodPost,
			wantFound: true,
			wantNode: &node{
				path:     "signUp",
				children: map[string]*node{},
				handler:  mockHandler,
			},
		},
		{
			name:      "GET 不存在handler",
			path:      "/user/profile",
			method:    http.MethodGet,
			wantFound: true,
			wantNode: &node{
				path:     "profile",
				children: map[string]*node{},
				handler:  mockHandler,
			},
		},
		{
			name:      "参数路径",
			path:      "/user/signUp/:id",
			method:    http.MethodGet,
			wantFound: true,
			wantNode: &node{
				path:    ":id",
				handler: mockHandler,
			},
		},
	}

	for _, tc := range testCases1 {
		t.Run(tc.name, func(t *testing.T) {
			n, ok := r.findRouter(tc.method, tc.path)
			if !ok {
				t.Log("未匹配路径")
				return
			}
			assert.Equal(t, ok, tc.wantFound)
			assert.Equal(t, n.path, tc.wantNode.path)
			msg, ok := n.equal(tc.wantNode)
			if !ok {
				t.Log(msg)
				return
			}
			nHandler := reflect.ValueOf(n.handler)
			wHandler := reflect.ValueOf(tc.wantNode.handler)
			assert.True(t, nHandler == wHandler)
		})
	}
}
