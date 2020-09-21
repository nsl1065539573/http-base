package gee

import (
	"log"
)
// 修改Handler为接收context 因为context中已经包含了 wr
type Handler func(*Context)

// 将路由抽出来作为一个结构体
type router struct {
	handlers map[string]Handler 
}

// 新建一个路由 应该在新建一个engine的时候新建
func newRouter() *router {
	return &router{
		make(map[string]Handler),
	}
}

// 新注册一个路由
func (r *router) addRouter(method, path string, handler Handler) {
	key := method + "-" + path
	log.Printf("Route: %s", key)
	r.handlers[key] = handler
}

// 处理http请求的方法
// 目前操作方式还是在注册路由时添加处理方法，在此处判断是否已经注册
// 并调用用户自己注册的方法
func (r *router) handle(c *Context) {
	key := c.r.Method + "-" + c.r.URL.Path
	log.Printf("get path: %s", key)
	log.Printf("router's routers %v", r.handlers)
	v, ok := r.handlers[key]
	log.Printf("did has key: %v", ok)
	if ok {
		v(c)
	} else {
		c.String(404 ,"%s not found aaa", c.r.URL.Path)
	}
}