package web_blaster

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Middleware 中间件的类型，此处采用函数式写法，使用中间件时，需传入这个类型的函数
// 这里会提供一些默认中间件，用户可以通过自身需要自己创建一些符合当前场景的中间件
// 如果使用Middleware，需要用户自己来控制是否需要调用下一个Middleware，
// 这里负责把用户的Middleware做成一条责任链，不断返回不断调用，但是否继续执行（比如在执行中间件时遇到一些不在预期内的错误），
// 需要用户自己来判断是否继续执行下一个Middleware
type Middleware func(HandlerFunc) HandlerFunc

// MiddlewareAccessLogBuilder 默认日志中间件
type MiddlewareAccessLogBuilder struct {
	AccessLogFunc func(accessLog []byte)
}

// Builder 打印日志
func (m *MiddlewareAccessLogBuilder) Builder() Middleware {
	m.AccessLogFunc = func(accessLog []byte) {
		log.Println(string(accessLog))
	}
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			body, _ := ioutil.ReadAll(c.req.Body)
			al := AccessLog{
				Method: c.req.Method,
				Path:   c.req.URL.String(),
				Route:  c.route,
				Body:   string(body),
			}
			bs, err := json.Marshal(al)
			if err == nil {
				m.AccessLogFunc(bs)
			}
			next(c)
			//log.Println(c.res)
		}
	}
}

// BodyNopCloser 将body改为可以重复读取的
func BodyNopCloser() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			c.req.Body = ioutil.NopCloser(c.req.Body)
			next(c)
		}
	}
}

// AccessLog 定义控制台打印日志的字段
type AccessLog struct {
	// 请求方法
	Method string
	// 请求路径
	Path string
	// 匹配上的路由
	Route string
	Body  string
}
