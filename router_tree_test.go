package server

import (
	"net/http"
	"testing"
)

func Test_RegisterRouter(_t *testing.T) {
	path := "user/hand/hah"
	pather := "role/hander"
	pathers := "role/handers"
	var f Hand
	f = func(http.ResponseWriter, *http.Request) {
		_t.Log("我被调用啦")
	}
	// 注册路由
	t := &tree{
		leafs: make(map[string]*leaf),
	}
	t.RegisterRouter(path, f, "get")
	t.RegisterRouter(pather, f, "get")
	t.RegisterRouter(pathers, f, "get")

	// 获取处理方法
	if f, err := t.GetHand(path, "get"); err != nil {
		_t.Error("注册/解析路由映射失败")
	} else {
		f(struct{ http.ResponseWriter }{}, &http.Request{})
		_t.Log("fullpath:", t)
	}
	_t.Log("tree:", t)
}
