package http

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"JMSpecialTrain/library/log"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/pkg/errors"
)

// Service http 服务端
type Service struct {
	Router

	// 请求方法不匹配时调用
	MethodNotAllow Hander

	// TODO: metastore is the path as key and the metadata of this path as value, it export via /metadata
	metastore map[string]map[string]interface{}
	mux       *http.ServeMux

	server atomic.Value  // store *http.Server
	stop   chan struct{} // 确保 server.Shutdown()执行完成再退出
}

// NewServer 新建 http 服务实例
func NewServer() *Service {
	s := &Service{
		Router: Router{
			root: true,
		},
		mux:       http.NewServeMux(),
		metastore: make(map[string]map[string]interface{}),
		stop:      make(chan struct{}),
	}
	s.Router.s = s
	return s
}

// RunHTTP 启动 http 服务
func (s *Service) RunHTTP(address string) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("http had closed:%+v", err)
		}
	}()
	server := &http.Server{
		Addr:         address,
		ReadTimeout:  time.Duration(30 * time.Second),
		WriteTimeout: time.Duration(30 * time.Second),
	}
	server.Handler = s.mux
	// TODO: 注册Shutdown方法(调用server.RegisterOnShutdown)解决long-live connection 和 hijacked connection
	s.server.Store(server)

	err := gracehttp.Serve(server)
	if err != nil {
		<-s.stop
		if err == http.ErrServerClosed {
			log.Warn("%+v:", err)
		}
	}
}

// addRouter 添加路由
func (s *Service) addRouter(method, path string, handlers ...Hander) {
	if path[0] != '/' {
		panic("path must begin with '/'")
	}
	if method == "" {
		panic("HTTP method can not be empty")
	}
	if _, ok := s.metastore[path]; !ok {
		s.metastore[path] = make(map[string]interface{})
	}
	s.metastore[path]["method"] = method
	s.mux.HandleFunc(path, func(resp http.ResponseWriter, req *http.Request) {
		c := &Context{
			Context:  nil,
			s:        s,
			Request:  req,
			Writer:   resp,
			method:   method,
			handlers: handlers,
		}

		ctx, cancel := context.WithCancel(context.Background())
		c.Context = ctx
		defer cancel()

		c.Next()
	})

	log.Info("注册路由:%s", path)
}

// Shutdown the http server without interrupting active connections.
func (s *Service) Shutdown(ctx context.Context) error {
	server := s.server.Load().(*http.Server)
	if server == nil {
		return errors.New("no server")
	}
	err := server.Shutdown(ctx)
	s.stop <- struct{}{}
	return err
}
