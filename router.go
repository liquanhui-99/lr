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
	// 校验请求路径是否合法
	r.validatePath(path)
	// 判断树是否存在，不存在则创建
	var root *node
	root, ok := r.trees[pattern]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[pattern] = root
	}

	// 根节点特殊处理
	if path == "/" {
		root.handler = handler
		return
	}

	// 去掉开头和结尾的/，切割path，并遍历每一个路径是否存在，不存在则创建
	path = strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			panic("请求路径不能包括连续的/")
		}
		children := root.childOperator(seg)
		root = children
	}
	root.handler = handler
}

func (r *router) validatePath(path string) {
	if path == "" {
		panic("路径不能为空")
	}

	if len(path) == 1 && path != "/" {
		panic("根路径不是/")
	}

	if len(path) > 1 {
		if path[0] != '/' {
			panic("请求路径不是以/开头")
		}
		if path[len(path)-1] == '/' {
			panic("请求路径不能以/结尾")
		}
	}
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
