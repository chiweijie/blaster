package web_blaster

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

// Context 服务的上下文，请求生命周期的控制
type Context struct {
	req       *http.Request
	res       http.ResponseWriter
	pathValue map[string]string
	route     string
}

// BindJSON 解析request中的body中的json到传入的结构体指针
func (c *Context) BindJSON(val any) error {
	if c.req.Body == nil {
		return errors.New("Body is nil")
	}
	decoder := json.NewDecoder(c.req.Body)
	err := decoder.Decode(val)
	if err == io.EOF {
		err = nil
	}
	return err
}

// FormValue 获取表单中的key对应的value
func (c *Context) FormValue(key string) value {
	err := c.req.ParseForm()
	if err != nil {
		return value{
			err: err,
		}
	}
	return value{
		val: c.req.FormValue(key),
	}
}

// FormValueOrDefault 获取表单中的key对应的value，如果表单中无数据，或者value为空，则返回默认值
func (c *Context) FormValueOrDefault(key string, def string) value {
	err := c.req.ParseForm()
	val := c.req.FormValue(key)
	if err != nil || val == "" {
		return value{
			val: def,
		}
	}
	return value{
		val: val,
	}
}

// PathValue 获取路径参数，如localhost:8080/detail:id 中的id
func (c *Context) PathValue(key string) value {
	val, ok := c.pathValue[key]
	if !ok {
		return value{
			err: errors.New("Not found key!"),
		}
	}
	return value{
		val: val,
	}
}

// WriteJson 写入json格式的响应
func (c *Context) WriteJson(code int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.res.WriteHeader(code)
	c.res.Write(bs)
	return err
}

// WriteString 写入string类型的响应
func (c *Context) WriteString(code int, val string) error {
	//c.res.WriteHeader(code)
	c.res.Write([]byte(val))
	return nil
}

// SetCookie 设置cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.res, cookie)
}

// value 返回解析后的具体数据或者error
type value struct {
	val string
	err error
}

// ToInt64 将获取的value转化为int64
func (v value) ToInt64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}
	return strconv.ParseInt(v.val, 10, 64)
}

// ToInt 将获取的value转化为int
func (v value) ToInt() (int, error) {
	if v.err != nil {
		return 0, v.err
	}
	return strconv.Atoi(v.val)
}

// ToUint64 将获取的value转化为uint64
func (v value) ToUint64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}
	return strconv.ParseInt(v.val, 10, 64)
}

// ToString 将获取的value转化为string
func (v value) ToString() (string, error) {
	if v.err != nil {
		return "", v.err
	}
	return v.val, v.err
}

// ToFloat64 将获取的value转化为float64
func (v *value) ToFloat64() (float64, error) {
	if v.err != nil {
		return 0, v.err
	}
	return strconv.ParseFloat(v.val, 64)
}

// Mes 写入数据
type Mes map[string]string
