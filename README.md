# blaster
web framework


// 实例代码

package test

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	blaster "web-blaster"
)

func Login(next blaster.HandlerFunc) blaster.HandlerFunc {
	return func(c *blaster.Context) {
		fmt.Println("用户是否登录")
		next(c)
	}
}

func Log(next blaster.HandlerFunc) blaster.HandlerFunc {
	return func(c *blaster.Context) {
		fmt.Println("log...")
		next(c)
	}
}

func TestWebTest(t *testing.T) {
	// 默认初始化服务，如需添加自己的 Middleware 需使用 blaster.NewHTTPServe 或 blaster.NewHTTPSServe，然后调用Use
	// 例如：
	//r := blaster.NewHTTPServe()
	//r.Use(Log, Login)
	r := blaster.DefaultHTTP()
	// 注册get路由
	r.Get("/", SomeLogic)
	// 注册带有通配符的路由
	r.Get("/star/*/abc", starNode)
	r.Get("/star1/*", starNode)
	//注册含有路径参数的路由
	r.Get("/pathValue/:id", paramNode)
	// 路由分组
	g := r.Group("/profix")
	g.Get("/test", func(c *blaster.Context) {
		c.WriteString(http.StatusOK, "hello group!")
	})
	// 路由分组
	g2 := r.Group("/api")
	u := new(User)
	{
		g2.Post("/json", func(c *blaster.Context) {
			err := c.BindJSON(u)
			if err != nil {
				fmt.Println(err)
			}
			c.WriteJson(http.StatusOK, blaster.Mes{
				"status": "ok",
			})
		})

		g2.Post("/form", func(c *blaster.Context) {
			age, _ := c.FormValue("age").ToInt()
			name, _ := c.FormValue("name").ToString()
			fmt.Printf("age type is %T, value is %v\n", age, age)
			fmt.Printf("name type is %T, value is %v\n", name, name)
			c.WriteJson(http.StatusOK, blaster.Mes{
				"status": "ok",
			})
		})
	}
	// 启动
	r.Start(":8080")
}

type User struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func SomeLogic(c *blaster.Context) {
	c.WriteString(http.StatusOK, "hello blaster")
}

func starNode(c *blaster.Context) {
	c.WriteString(http.StatusOK, "hello star node")
}

func paramNode(c *blaster.Context) {
	c.WriteString(http.StatusOK, "hello pathValue node")
	val, err := c.PathValue("id").ToString()
	if err != nil {
		log.Println(err)
		return
	}
	c.WriteString(http.StatusOK, val)
}
