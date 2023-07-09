package lorm

import "strings"

// router 一个森林，每一种请求方式都是一颗树，
// key是请求的方法，value是树的节点
type router struct {
	trees map[string]*node
}

type node struct {
	// 请求路径
	path string
	// 子节点， key是路径
	children map[string]*node
	// 业务处理逻辑
	handler HandleFunc
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

// AddRouter 注册路由信息
func (r *router) AddRouter(pattern, path string, handler HandleFunc) {
	root, ok := r.trees[pattern]
	if !ok {
		root = &node{
			path:     "/",
			children: make(map[string]*node),
		}
		r.trees[pattern] = root
	}
	path = strings.TrimLeft(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		// 递归查找每一个子节点是否存在
		children := root.childOrCreate(seg)
		root = children
	}
	root.handler = handler
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
