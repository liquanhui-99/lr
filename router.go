package lorm

import "strings"

// router the trees of http router. map key is http method.
type router struct {
	trees map[string]*node
}

// node the router's tree.
type node struct {
	path     string
	children map[string]*node
	handler  HandleFunc
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

func (r *router) AddRouter(method, path string, handler HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	path = path[1:]
	for _, seg := range strings.Split(path, "/") {
		children := root.childOrCreate(seg)
		root = children
	}
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}
