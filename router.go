package lr

import (
	"fmt"
	"strings"
)

type HandleFunc func(Context)

// router 路由森林，不是单颗树，key是请求的方法，value是单个树，树上有各个路由节点
type router struct {
	trees map[string]*node
}

type node struct {
	// 请求的路径
	path string
	// 自路由，key是下一个子路径(例如：路径/user/signIn，user是当前的path，signIn是children中的key)
	// value是node的节点，每一个路径都会有自己的节点信息
	children map[string]*node
	// 处理具体的业务逻辑
	handler HandleFunc
}

func NewRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// addRouter 添加路由 先查看路由树中是否存在，不存在创建路由节点，查询前需要先对添加的路径做特殊校验
// 1. 必须以 / 开头
// 2. 不能以 / 结尾
// 3. 不能是空字符串
// 4. 不能是连续的 ///，无论是开头、结尾、还是路径中间
func (r *router) addRouter(method, path string, handler HandleFunc) {
	if len(path) == 0 {
		panic("请求路径不能为空")
	}

	root, ok := r.trees[method]
	if !ok {
		// 根节点不存在，需要先创建根节点
		root = &node{
			path:     "/",
			children: map[string]*node{},
		}
		r.trees[method] = root
	}

	if path[0] != '/' {
		panic("请求路径必须以/开头")
	}

	if path != "/" && path[len(path)-1] == '/' {
		panic("请求路径不能以/结尾")
	}

	// 处理请求路径是根路径
	if path == "/" {
		// 需要单独处理根节点冲突的问题
		if root.handler != nil {
			panic("路由冲突，重复注册[/]")
		}
		root.handler = handler
		return
	}

	segments := strings.Split(path[1:], "/")
	for _, seg := range segments {
		if seg == "" {
			panic("请求路径不能包含连续的/")
		}
		child := root.childOf(seg)
		root = child
	}

	// 这里处理普通路径重复注册的问题
	if root.handler != nil {
		panic(fmt.Sprintf("路由冲突，重复注册[%s]", path))
	}

	root.handler = handler
}

// findRouter 匹配路由
func (r *router) findRouter(method, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	// 根节点需要单独处理
	if path == "/" {
		return root, true
	}

	path = strings.Trim(path, "/")
	segments := strings.Split(path, "/")
	for _, seg := range segments {
		child, ok := root.matchChildOf(seg)
		if !ok {
			return nil, false
		}
		root = child
	}

	if root.handler == nil {
		return nil, false
	}

	// 返回节点和true，调用者知道有这个节点，但是节点的handler是不是目标handler需要自己判断
	return root, true
}

func (n *node) matchChildOf(seg string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}

	child, ok := n.children[seg]
	return child, ok
}

func (n *node) childOf(seg string) *node {
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	nd, ok := n.children[seg]
	if !ok {
		nd = &node{
			path:     seg,
			children: map[string]*node{},
		}
		n.children[seg] = nd
	}

	return nd
}
