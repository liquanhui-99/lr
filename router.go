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
	// 判断树是否存在，不存在则创建
	var root *node
	root, ok := r.trees[pattern]
	if !ok {
		root = &node{
			path:     "/",
			children: map[string]*node{},
		}
		r.trees[pattern] = root
	}

	// 切割path，并遍历每一个路径是否存在，不存在则创建
	path = strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		children := root.childOperator(seg)
		root = children
	}
	root.handler = handler
}

func (n *node) childOperator(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		res = &node{
			path:     seg,
			children: map[string]*node{},
		}
		n.children[seg] = res
	}
	return res
}
