package gee

import (
	"net/http"
	"fmt"
	"encoding/json"
)

/**
* Context 结构体  用来封装上下文
* 目前包含 writer request path method statusCode等字段
* writer 用于写响应
* request 用于接收请求的数据
* Path 记录路径
* Method 请求的方法
* StatusCode http状态码
*/

type J map[string]interface{}

type Context struct {
	w http.ResponseWriter
	r *http.Request
	Path string
	Method string
	StatusCode int
}

// 新建一个context， 需要传入 w以及r
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w: w,
		r: r,
		Path: r.URL.Path,
		Method: r.Method,
	}
}

// 查询form表单的数据，根据key查询value
func (c *Context) PostForm(key string) string {
	return c.r.FormValue(key)
}

// 拿到URL中的参数值， 比如 http://localhost:9999?name=Tom  c.Query("name") 返回 "Tom" 
func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

// 设置http状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.w.WriteHeader(code)
}

// 设置http返回头信息
func (c *Context) setHeader(key, value string) {
	c.w.Header().Set(key, value)
}

// 返回String
func (c *Context) String(code int, format string, args ...interface{}) {
	c.setHeader("Content-Type", "text/plain")
	c.Status(code)
	fmt.Fprintf(c.w, format, args...)
}

// 返回JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.setHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.w, err.Error(), 500)
	}
}

// 返回数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.w.Write(data)
}

// 返回HTML文本
func (c *Context) HTML(code int, html string) {
	c.setHeader("Content-Type", "text/html")
	c.Status(code)
	c.w.Write([]byte(html))
}