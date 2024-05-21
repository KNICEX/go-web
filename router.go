package web

import "strings"

type router struct {

	// method -> tree root
	trees map[string]*node
}

type node struct {
	path string
	// 注册的路由字符串
	route    string
	children []*node
	// 通配符
	starChild  *node
	paramChild *node
	handlers   []HandleFunc
}

type matchInfo struct {
	node       *node
	pathParams map[string]string
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

// addRoute 注册路由
// - 不能同时注册多个相同的路由
// - 不能在同一个位置同时有通配符和路径参数
func (r *router) addRoute(method string, path string, handlers ...HandleFunc) {
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			continue
		}

		child := root.childOrCreate(seg)

		root = child
	}

	if root.handlers != nil {
		panic("duplicated path")
	}
	root.route = path
	root.handlers = append(root.handlers, handlers...)
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	segs := strings.Split(path, "/")
	pathParams := make(map[string]string, 2)
	for _, seg := range segs {
		if seg == "" {
			continue
		}

		child, isParam, ok := root.childOf(seg)
		if !ok {
			return nil, false
		}
		if isParam {
			pathParams[child.path[1:]] = seg
		}

		root = child

	}

	return &matchInfo{
		node:       root,
		pathParams: pathParams,
	}, true
}

func (n *node) childOrCreate(seg string) *node {
	if seg == "*" {
		if n.paramChild != nil {
			panic("can't register wildcard and param node at the same time")
		}
	}

	if seg[0] == ':' {
		if n.starChild != nil {
			panic("can't register wildcard and param node at the same time")
		}
		if n.paramChild != nil {
			panic("can't register two param nodes at the same level")
		}
	}

	child, _, ok := n.childOf(seg)
	if !ok {
		child = &node{
			path: seg,
		}
		if seg == "*" {
			n.starChild = child
		} else if seg[0] == ':' {
			n.paramChild = child
		} else {
			n.children = append(n.children, child)
		}
	} else if child.path == "*" {
		// 匹配到通配符
		if seg == "*" {
			// 自己也是通配符，不需要创建新的节点
			return child
		} else {
			// 自己是详细路径，创建一个新的节点
			child = &node{
				path: seg,
			}
			n.children = append(n.children, child)
		}
	} else if child.path[0] == ':' {
		// 已有参数节点，自己是详细路径，创建一个新的节点
		child = &node{
			path: seg,
		}
		n.children = append(n.children, child)
	}
	return child
}

// childOf returns the child node of n with path seg
// 第一个返回值是找到的节点
// 第二个返回值是找到的节点是否是路径参数节点
// 第三个返回值是否找到了节点
func (n *node) childOf(seg string) (*node, bool, bool) {
	if n.children == nil {

		if n.paramChild != nil {
			return n.paramChild, true, true
		}

		return n.starChild, false, n.starChild != nil
	}

	for _, child := range n.children {
		if child.path == seg {
			return child, false, true
		}
	}

	// 静态节点中没有找到，查找通配符节点
	if n.paramChild != nil {
		return n.paramChild, true, true
	}
	return n.starChild, false, n.starChild != nil
}
