package http

// Hander 处理方法
type Hander func(*Context)

// Router 路由管理
type Router struct {
	handers  []Hander
	basePath string
	root     bool
	s        *Service
}

func (r *Router) add(method, path string, handlers ...Hander) {
	path = joinPaths(r.basePath, path)
	if r.handers == nil {
		r.handers = make([]Hander, 0, 5)
		r.handers = append(r.handers, MethodNotAllow)
	}
	for _, h := range handlers {
		r.handers = append(r.handers, h)
	}
	r.s.addRouter(method, path, r.handers...)
}

// Group 创建新的 Router, 使用此 Router 生成的路由都将拥有前缀: r.basePath+basePath, 并且执行 handers
func (r *Router) Group(basePath string, handers ...Hander) *Router {
	return &Router{
		handers:  handers,
		basePath: joinPaths(r.basePath, basePath),
		s:        r.s,
		root:     false,
	}
}

// GET 注册 GET 请求的路由
func (r *Router) GET(path string, handlers ...Hander) {
	r.add("GET", path, handlers...)
}

// POST 注册 POST 请求的路由
func (r *Router) POST(path string, handlers ...Hander) {
	r.add("POST", path, handlers...)
}

// PUT 注册 PUT 请求的路由
func (r *Router) PUT(path string, handlers ...Hander) {
	r.add("PUT", path, handlers...)
}

// DELETE 注册 DELETE 请求的路由
func (r *Router) DELETE(path string, handlers ...Hander) {
	r.add("DELETE", path, handlers...)
}
