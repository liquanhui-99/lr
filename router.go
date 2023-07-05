package lorm

// Router 每一种都是一颗树，key是请求的方法，value是树的节点
type Router struct {
	router map[string]*node
}

type node struct {
	path     string  // 请求路径
	children []*node // 子节点，key是
	handler  HandleFunc
}
