package lr

import (
	"fmt"
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

	var mockHanlerFunc HandleFunc = func(ctx Context) {}
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

	r := NewRouter()
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
