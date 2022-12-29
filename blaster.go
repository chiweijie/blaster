package web_blaster

import (
	"net/http"
	"os"
)

// Serve 核心接口，定义服务的核心方法
// 这个接口只保留了服务最核心的方法，之前参考过其他框架的核心接口设计，
// 大多都是保留了十多个方法，而很多方法都是基于已有方法的封装，
// 而且这样设计的大多数情况都是侵入核心了，
// 个人看法没有什么必要，所以这里采用非侵入式的做法，只保留最核心的方法，
type Serve interface {
	http.Handler

	// 注册路由，此方法不对外暴露，用户使用的话应该使用该方法的封装
	addRoute(method string, path string, handleFunc HandlerFunc)

	// 启动函数
	Start(string) error
}

// HTTPServe http服务
type HTTPServe struct {
	// 路由树
	*router

	// 保存注册的中间件
	mw []Middleware
}

// HTTPSServe https服务
type HTTPSServe struct {
	*HTTPServe
	certFile string
	keyFile  string
}

// ServeHTTP 此方法实现了http.Handler接口，接受到请求后会执行此方法
func (h *HTTPServe) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		res: res,
		req: req,
	}
	root := h.serve
	// 这里之所以采用这种写法，是因为我们需要用下一个方法来构造上一个方法，最后一层一层不断返回，调用执行
	for i := len(h.mw) - 1; i >= 0; i-- {
		root = h.mw[i](root)
	}
	root(ctx)
}

// Use 添加中间件
// Use函数正常来说只能调用一次，如果发生再次调用，会将先前的Middleware覆盖
func (h *HTTPServe) Use(mw ...Middleware) {
	h.mw = mw
}

// Start http启动函数
func (h *HTTPServe) Start(addr string) error {
	return http.ListenAndServe(addr, h)
}

// serve 先查找路由，找到的话执行路由对应的方法，没找到的话返回404
func (h *HTTPServe) serve(ctx *Context) {
	// 查找路由，如果未找到，或者找到的节点是中间节点
	// 比如/user/detail/account，而请求的路径是/user/detail,
	// 则这里的ok会返回false，此时返回404错误
	mi, ok := h.findRoute(ctx.req.Method, ctx.req.URL.String())
	if !ok {
		ctx.res.WriteHeader(http.StatusNotFound)
		ctx.res.Write([]byte("404 PAGE NOT FOUND!"))
		return
	}
	ctx.pathValue = mi.param
	ctx.route = mi.n.path
	mi.n.handler(ctx)
}

// Start https启动函数
func (h *HTTPSServe) Start(addr string) error {
	return http.ListenAndServeTLS(addr, h.certFile, h.keyFile, h)
}

// Get 执行get请求，以下方法类似
func (h *HTTPServe) Get(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

func (h *HTTPServe) Head(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodHead, path, handleFunc)
}

func (h *HTTPServe) Post(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodPost, path, handleFunc)
}

func (h *HTTPServe) Put(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodPut, path, handleFunc)
}

func (h *HTTPServe) Patch(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodPatch, path, handleFunc)
}

func (h *HTTPServe) Delete(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodDelete, path, handleFunc)
}

func (h *HTTPServe) Connect(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodConnect, path, handleFunc)
}

func (h *HTTPServe) Options(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodOptions, path, handleFunc)
}

func (h *HTTPServe) Trace(path string, handleFunc HandlerFunc) {
	h.addRoute(http.MethodTrace, path, handleFunc)
}

// NewHTTPServe 初始化Serve
func NewHTTPServe() *HTTPServe {
	return &HTTPServe{
		newRoute(),
		nil,
	}
}

// DefaultHTTP 默认使用打印到终端的日志中间件，如果需要输出到文件，
// 需要使用NewHTTPServe()函数，然后使用Use方法自由添加中间件。
func DefaultHTTP() *HTTPServe {
	s := NewHTTPServe()
	m := &MiddlewareAccessLogBuilder{}
	s.Use(BodyNopCloser(), m.Builder())
	return s
}

// NewHTTPSServe 初始化Serve
func NewHTTPSServe(certFile string, keyFile string) *HTTPSServe {
	return &HTTPSServe{
		HTTPServe: NewHTTPServe(),
		certFile:  certFile,
		keyFile:   keyFile,
	}
}

// DefaultHTTPS 默认使用打印到终端的日志中间件，如果需要输出到文件，
// 需要使用 DefaultHTTPS() 函数，然后使用 Use 方法自由添加中间件。
func DefaultHTTPS(certFile string, keyFile string) *HTTPSServe {
	s := NewHTTPSServe(certFile, keyFile)
	m := &MiddlewareAccessLogBuilder{}
	s.Use(BodyNopCloser(), m.Builder())
	return s
}

type HandlerFunc func(*Context)

// Shutdown 方法
func (h *HTTPServe) Shutdown() {
	os.Exit(100)
}
