package web_blaster

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_addRoute(t *testing.T) {
	var handlerFunc HandlerFunc = func(context *Context) {}
	tests := []struct {
		name        string
		method      string
		path        string
		handlerFunc HandlerFunc
	}{
		{
			name:        "/",
			method:      http.MethodGet,
			path:        "/",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/user",
			method:      http.MethodGet,
			path:        "/user",
			handlerFunc: handlerFunc,
		},
		{
			name:        "user",
			method:      http.MethodGet,
			path:        "user",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/home/",
			method:      http.MethodGet,
			path:        "/home/",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/home/*/detail",
			method:      http.MethodGet,
			path:        "/home/*/detail",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/home/:order_id",
			method:      http.MethodGet,
			path:        "/home/:order_id",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/user/detail/profile",
			method:      http.MethodGet,
			path:        "/user/detail/profile",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/user/*",
			method:      http.MethodGet,
			path:        "/user/*",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/user/detail/profile",
			method:      http.MethodPost,
			path:        "/user/detail/profile",
			handlerFunc: handlerFunc,
		},
		{
			name:        "/business/order/detail/part",
			method:      http.MethodGet,
			path:        "/business/order/detail/part",
			handlerFunc: handlerFunc,
		},
	}
	wantRoute := &router{
		map[string]*node{
			"GET": &node{
				path:    "/",
				handler: handlerFunc,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: handlerFunc,
						starChild: &node{
							path:    "*",
							handler: handlerFunc,
						},
						children: map[string]*node{
							"detail": &node{
								path: "detail",
								children: map[string]*node{
									"profile": &node{
										path:    "profile",
										handler: handlerFunc,
									},
								},
							},
						},
					},
					"home": &node{
						path: "home",
						starChild: &node{
							path: "*",
							children: map[string]*node{
								"detail": &node{
									path:    "detail",
									handler: handlerFunc,
								},
							},
						},
						paramChild: &node{
							path:    ":order_id",
							handler: handlerFunc,
						},
						handler: handlerFunc,
					},
					"business": &node{
						path: "business",
						children: map[string]*node{
							"order": &node{
								path: "order",
								children: map[string]*node{
									"detail": &node{
										path: "detail",
										children: map[string]*node{
											"part": &node{
												path:    "part",
												handler: handlerFunc,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"POST": &node{
				path: "/",
				children: map[string]*node{
					"user": &node{
						path: "user",
						children: map[string]*node{
							"detail": &node{
								path: "detail",
								children: map[string]*node{
									"profile": &node{
										path:    "profile",
										handler: handlerFunc,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	r := &router{
		map[string]*node{},
	}

	for _, tt := range tests {
		r.addRoute(tt.method, tt.path, tt.handlerFunc)
	}
	errStr, ok := wantRoute.equal(r)
	assert.True(t, ok, errStr)

	findRoutes := []struct {
		name   string
		method string
		path   string

		wantPath  string
		wantFound bool
	}{
		{
			name:   "/",
			method: http.MethodGet,
			path:   "/",

			wantPath:  "/",
			wantFound: true,
		},
		{
			name:   "/user",
			method: http.MethodGet,
			path:   "/user",

			wantPath:  "user",
			wantFound: true,
		},
		{
			name:   "/user/detail/profile",
			method: http.MethodGet,
			path:   "/user/detail/profile",

			wantPath:  "profile",
			wantFound: true,
		},
		{
			name:   "/user/abc",
			method: http.MethodGet,
			path:   "/user/abc",

			wantPath:  "*",
			wantFound: true,
		},
		{
			name:   "/user/abc",
			method: http.MethodGet,
			path:   "/user/abc",

			wantPath:  "*",
			wantFound: true,
		},
		{
			name:   "/home/abc/detail",
			method: http.MethodGet,
			path:   "/home/abc/detail",

			wantPath:  "detail",
			wantFound: true,
		},
		{
			name:   "/user/detail",
			method: http.MethodGet,
			path:   "/user/detail",

			wantPath:  "detail",
			wantFound: false,
		},
	}
	for _, f := range findRoutes {
		t.Run(f.name, func(tt *testing.T) {
			node, found := r.findRoute(f.method, f.path)
			assert.Equal(tt, found, f.wantFound)
			if !found {
				return
			}
			assert.Equal(tt, node.n.path, f.wantPath)
			//assert.NotNil(tt, )
		})
	}
}

func (r router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.handler)
	yhv := reflect.ValueOf(y.handler)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}
