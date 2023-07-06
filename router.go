package lorm

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

}
