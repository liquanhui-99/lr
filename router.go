package lorm

import (
	"fmt"
	"strings"
)

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
	// 请求路径上的参数(/user/login/:123)
	paramChild *node
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

// addRouter 注册路由信息
func (r *router) addRouter(pattern, path string, handler HandleFunc) {
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
		if root.handler != nil {
			panic(fmt.Sprintf("请求路径[%s]已注册", path))
		}
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
	if root.handler != nil {
		panic(fmt.Sprintf("请求路径[%s]已注册", path))
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
	// 处理请求路径上的参数
	if seg[0] == ':' {
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	// 静态路由处理
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

// match 匹配请求路径，优先匹配静态路由，再匹配路径参数
// *node 匹配的节点信息
// bool 是否是路径参数
// bool 是否匹配到该路径
func (n *node) match(seg string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
	}

	child, ok := n.children[seg]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return nil, false, n.paramChild != nil
	}
	return child, false, ok
}

func (r *router) matchRouter(pattern, path string) (*pathInfo, bool) {
	tree, ok := r.trees[pattern]
	if !ok {
		return nil, false
	}

	path = strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		children, pt, ok := tree.match(seg)
		if !ok {
			return nil, false
		}
		// 命中了路径参数
		if pt {
			return &pathInfo{
				n:          children,
				pathParams: map[string]string{},
			}, true
		}
		tree = children
	}
	// 判断节点存在的情况下，是否有handler
	return &pathInfo{}, tree.handler != nil
}

// pathInfo 路径信息
type pathInfo struct {
	// 节点信息
	n *node
	// 路径参数，key是参数名，value是参数对应的值
	pathParams map[string]string
}
