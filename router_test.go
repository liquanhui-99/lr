package lorm

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter(t *testing.T) {
	testCases := []struct {
		path    string
		method  string
		handler HandleFunc
	}{
		{
			path:   "/root/user",
			method: "GET",
		},
	}
	rt := newRouter()
	mockHandler := func(ctx Context) {}
	for _, tc := range testCases {
		rt.AddRouter(tc.method, tc.path, tc.handler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "root",
				handler: nil,
				children: map[string]*node{
					"/user": &node{
						path:     "/user",
						children: map[string]*node{},
						handler:  mockHandler,
					},
				},
			},
		},
	}

	msg, equal := rt.equal(wantRouter)
	if !equal {
		fmt.Println(msg)
		return
	}
	fmt.Println("success!")
}

func (r *router) equal(dest *router) (string, bool) {
	for key, value := range r.trees {
		destValue, ok := dest.trees[key]
		if !ok {
			return "http method not exist!", false
		}
		msg, equal := value.equal(destValue)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(dest *node) (string, bool) {
	if n.path != dest.path {
		return "router not exist!", false
	}

	nHandler := reflect.ValueOf(n.handler)
	destHandler := reflect.ValueOf(dest.handler)
	if nHandler != destHandler {
		return "handler not equal!", false
	}

	if len(n.children) != len(dest.children) {
		return "children error", false
	}

	for path, value := range n.children {
		dst, ok := dest.children[path]
		if !ok {
			return "child node not exist!", false
		}
		msg, equal := value.equal(dst)
		if !equal {
			return msg, false
		}
	}

	return "", true
}
