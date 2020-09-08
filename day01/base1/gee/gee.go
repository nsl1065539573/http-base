package gee

import (
	"fmt"
	"net/http"
)

type Handler func(w http.ResponseWriter, req *http.Request)

// 继承 http.Handler
type engine struct {
	router map[string]Handler
}

// 实现 http.Handler   判断URL与method是否已经注册
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if v, ok := e.router[key]; ok {
		v(w, req)
	} else {
		fmt.Fprintf(w, "404 not found\n")
	}
}

// 新建一个engine
func New() *engine {
	return &engine{make(map[string]Handler)}
}
 
// 给engine添加一个处理器
func (e *engine) addRouter(path string, handler Handler) {
	e.router[path] = handler
}

// 声明一个GET方法的处理器
func (e *engine) GET(path string, handler Handler) {
	key := "GET-" + path
	e.addRouter(key, handler)
}

// 声明一个POST方法的处理器
func (e *engine) POST(path string, handler Handler) {
	key := "POST-" + path
	e.addRouter(key, handler)
}

// 运行server
func (e *engine) Run(port string) (err error) {
	return http.ListenAndServe(port, e)
}
