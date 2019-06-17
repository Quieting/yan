package http

var (
	// ErrMethodNotAllow 请求方法不匹配时返回的错误信息
	ErrMethodNotAllow = "invalid request"
)

// MethodNotAllow 路由注册时会默认添加此方法, 可覆盖
var MethodNotAllow = Hander(func(c *Context) {
	if c.method != c.Request.Method {
		c.Writer.Write([]byte(ErrMethodNotAllow))
		c.Abort()
	}
})
