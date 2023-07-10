package lorm

import (
	"fmt"
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
			name:    "path",
			pattern: "POST",
			path:    "/api/user/login",
		},
	}
	var mockHandler HandleFunc = func(c Context) {}
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
		},
	}
	router := newRouter()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router.AddRouter(tc.pattern, tc.path, mockHandler)
		})
	}
	msg, ok := router.equal(wantRouter)
	if !ok {
		fmt.Println("结果为: ", msg)
		return
	}
	fmt.Println("匹配成功！")
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
