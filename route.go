package web_blaster

import (
	"net/http"
	"strings"
)

// router 路由树
type router struct {
	// 每一类方法对应一棵树，
	// 比如如果注册了GET，POST，DELETE方法，那么就是三颗树
	// 因为这里的路由树是用map实现的，所以并不能支持动态并发添加路由，项目一旦启动，则不能再添加路由
	// 如果需要动态添加路由，可以采用装饰器模式封装一下这个结构体
	// type SafeRouter struct {
	// 		l sync.RWMutex
	//  	router
	// }
	// 当然目前这个结构体是私有的，在后期会慢慢补充动态并发操作的api
	trees map[string]*node
}

// node 路由树的节点
type node struct {
	// 请求路径
	path string
	// 子节点
	children map[string]*node
	// 通配符节点
	starChild *node
	// 路径参数节点
	paramChild *node
	// 该路径对应的方法
	handler HandlerFunc
	// 匹配上的路由
	Route string
}

// matchInfo 查找到路由后返回的路由节点及参数
type matchInfo struct {
	// 查找到的路由返回最终节点
	n *node
	// 保存该节点对应的路径参数
	param map[string]string
}

// Group 路由分组结构体
type Group struct {
	// 分组的路由前缀
	prefix string
	s      Serve
}

// Group 分组路由，保存传入的前缀到Group实例中
func (hs *HTTPServe) Group(prefix string) *Group {
	return &Group{
		prefix,
		hs,
	}
}

// addRoute 分组的基础上添加路由
func (g *Group) addRoute(method string, path string, handlerFunc HandlerFunc) {
	g.s.addRoute(method, g.prefix+path, handlerFunc)
}

// Get 执行get请求，以下方法类似
// 方法待改进，使用上这里没有问题，但是跟HTTPServe结构体中的方法耦合了，后面会对这里进行改进
func (g *Group) Get(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodGet, path, handleFunc)
}

func (g *Group) Head(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodHead, path, handleFunc)
}

func (g *Group) Post(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodPost, path, handleFunc)
}

func (g *Group) Put(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodPut, path, handleFunc)
}

func (g *Group) Patch(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodPatch, path, handleFunc)
}

func (g *Group) Delete(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodDelete, path, handleFunc)
}

func (g *Group) Connect(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodConnect, path, handleFunc)
}

func (g *Group) Options(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodOptions, path, handleFunc)
}

func (g *Group) Trace(path string, handleFunc HandlerFunc) {
	g.addRoute(http.MethodTrace, path, handleFunc)
}

// newRoute 初始话一个路由，暂时不指定任何参数，后续会新增函数或改造此函数
func newRoute() *router {
	return &router{
		map[string]*node{},
	}
}

// addRoute 注册路由
func (r *router) addRoute(method string, path string, handlerFunc HandlerFunc) {
	tree, ok := r.trees[method]
	// 如果还没有此方法的路由树，新建一颗
	if !ok {
		tree = &node{
			path: "/",
		}
		r.trees[method] = tree
	}
	// 将注册的方法赋值给对应路由
	defer func() {
		tree.handler = handlerFunc
	}()
	// 如果注册的路由为“/”，无需后续处理
	if path == "/" {
		return
	}
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")

	// 构建路由树
	for _, s := range segs {
		child := tree.getChild(s)
		tree = child
	}
	tree.Route = path
}

// getChild 获取子节点，若不存在则新建子节点
func (n *node) getChild(path string) *node {
	// 如果当前路径是“*”，那么把他添加到当前节点的starChild中
	if path == "*" {
		if n.starChild == nil {
			n.starChild = &node{
				path: path,
			}
		}
		return n.starChild
	}
	// 如果当前路由是由“:”开头的，那么把他放到当前节点的paramChild中
	if path[0] == ':' {
		if n.paramChild == nil {
			n.paramChild = &node{
				path: path,
			}
		}
		return n.paramChild
	}

	if n.children == nil {
		n.children = make(map[string]*node, 1)
	}

	child, ok := n.children[path]

	if !ok {
		child = &node{path: path}
		n.children[path] = child
	}
	return child
}

// findRoute 根据请求查找路由
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	tree, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{
			n: tree,
		}, true
	}

	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		// 找不到路径直接返回nil和false，否则将tree的指针指向找到的子节点
		child, ok := tree.children[seg]
		if !ok {
			if seg[0] == ':' && tree.paramChild != nil {
				return &matchInfo{
					n: tree.paramChild,
					param: map[string]string{
						tree.paramChild.path[1:]: seg[1:],
					},
				}, true
			}
			if tree.starChild != nil {
				tree = tree.starChild
				continue
			}
			return nil, false
		}
		tree = child
	}
	// 此处如果handler为nil，那么说明该路径是个中间节点，并未实际注册
	return &matchInfo{
		n: tree,
	}, tree.handler != nil
}
