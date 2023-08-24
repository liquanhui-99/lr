package lr

import "strings"

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

func (r *router) addRouter(method, path string, handler HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		// 根节点不存在，需要先创建根节点
		root = &node{
			path:     "/",
			children: map[string]*node{},
		}
		r.trees[method] = root
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	for _, seg := range segs {
		children := root.childOf(seg)
		root = children
	}
	root.handler = handler
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
