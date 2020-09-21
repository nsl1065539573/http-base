package gee

import (
	"strings"
	"net/http"
	"fmt"
)

type Handler func(*Context)

type router struct {
	roots map[string]*TrieNode
	handlers map[string]Handler
}

// 新建路由
func newRouter() *router {
	return &router{
		roots: make(map[string]*TrieNode),
		handlers: make(map[string]Handler),
	}
}

// 格式化路由
// 将string的path格式化为[]string的parts
func parseParts(path string) []string {
	parts := strings.Split(path, "/")
	res := make([]string, 0, len(parts))
	for i := 0; i < len(parts); i++ {
		item := parts[i]
		if item != "" {
			res = append(res, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return res
}

// 添加路由
func (r *router) addRouter(method, path string, handler Handler) {
	parts := parseParts(path)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = NewTrie()
	}
	r.roots[method].Insert(parts, path)
	
	key := method + "-" + path
	r.handlers[key] = handler
}

// 获取路由
func (r *router) getRouter(method, path string) (*TrieNode, map[string]string) {
	searchParts := parseParts(path)
	n, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	root := n.Search(searchParts)
	res := make(map[string]string)
	if root != nil {
		parts := strings.Split(root.path, "/")
		fmt.Printf("打印一下path::%v以及parts%v", root.path, parts)
		// 因为split会截出多一个空字符串在第一个
		if len(parts) > 0 && parts[0] == "" {
			parts = parts[1:]
		}
		for i := 0; i < len(parts); i++ {
			item := parts[i]
			if item[0] == ':' {
				res[item[1:]] = searchParts[i]
			}
			if item[0] == '*' && len(item) > 1 {
				res[item[1:]] = strings.Join(searchParts[i :], "/")
				break
			}
		}
		return root, res
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.path
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}