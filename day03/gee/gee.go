package gee

import (
	"net/http"
)

// 继承 http.Handler
type engine struct {
	router *router
}

// 实现 http.Handler  处理请求的方法在route结构体控制
func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	e.router.handle(NewContext(w, req))
}

// 新建一个engine
func New() *engine {
	return &engine{newRouter()}
}
 
// 给engine添加一个路由
func (e *engine) addRouter(method, path string, handler Handler) {
	e.router.addRouter(method, path, handler)
}

// 声明一个GET方法的路由
func (e *engine) GET(path string, handler Handler) {
	e.addRouter("GET", path, handler)
}

// 声明一个POST方法的路由
func (e *engine) POST(path string, handler Handler) {
	e.addRouter("POST", path, handler)
}

// 运行server
func (e *engine) Run(port string) (err error) {
	return http.ListenAndServe(port, e)
}
