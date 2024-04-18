package web

import "strings"

type router struct {

	// method -> tree root
	trees map[string]*node
}

type node struct {
	path     string
	children []*node
	// 通配符
	startChild *node
	handlers   []HandleFunc
}

func newRouter() *router {
	return &router{
		trees: make(map[string]*node),
	}
}

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
		child, ok := root.childOf(seg)
		if !ok {
			child = &node{
				path: seg,
			}
			if seg == "*" {
				root.startChild = child
			} else {
				root.children = append(root.children, child)
			}
		} else if child.path == "*" {
			// 匹配到通配符
			if seg == "*" {
				// 自己也是通配符，不需要创建新的节点
				child = child.startChild
			} else {
				// 自己是详细路径，创建一个新的节点
				child = &node{
					path: seg,
				}
				root.children = append(root.children, child)
			}

		}

		root = child
	}

	if root.handlers != nil {
		panic("duplicated path")
	}
	root.handlers = append(root.handlers, handlers...)
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			continue
		}

		child, ok := root.childOf(seg)
		if !ok {
			return nil, false
		}

		root = child

	}

	return root, true
}

// childOf returns the child node of n with path seg
func (n *node) childOf(seg string) (*node, bool) {
	if n.children == nil {
		return n.startChild, n.startChild != nil
	}

	for _, child := range n.children {
		if child.path == seg {
			return child, true
		}
	}

	// 静态节点中没有找到，查找通配符节点
	return n.startChild, n.startChild != nil
}
