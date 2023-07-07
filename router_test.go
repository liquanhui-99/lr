package lorm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRouter_AddRouter(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		pattern string
	}{
		{},
	}

	var mockHandler HandleFunc = func(ctx Context) {}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := newRouter()
			router.AddRouter(tc.pattern, tc.path, mockHandler)
		})
	}

}

func (r *router) equal(d *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := d.trees[k]
		if !ok {
			return fmt.Sprintf("请求方法不匹配"), false
		}
		// 比较每一个节点树是否相等
		msg, ok := v.equal(dst)
		if !ok {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(d *node) (string, bool) {
	// 比较当前节点的路径
	if n.path != d.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}

	// 比较handler
	nHandler := reflect.ValueOf(n.handler)
	dHandler := reflect.ValueOf(d.handler)
	if nHandler != dHandler {
		return fmt.Sprintf("HandlerFunc不想等"), false
	}

	// 比较子节点的路径数量
	if len(n.children) != len(d.children) {
		return fmt.Sprintf("子节点路径数量不一致"), false
	}

	for k, v := range n.children {
		val, ok := d.children[k]
		if !ok {
			return fmt.Sprintf("子节点路径不匹配"), false
		}
		msg, ok := v.equal(val)
		if !ok {
			return msg, false
		}
	}

	return "", true
}
