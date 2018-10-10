package server

import (
	"fmt"
	"http/util"
	"net/http"
)

func Start() {
	// 注册路由
	t.RegisterRoulter("auth", AuthorizationDigest, "get")

	// 启动服务
	s := http.NewServeMux()
	s.HandleFunc("/", hand)
	err := http.ListenAndServe(":7766", s)
	if err != nil {
		panic(err)
	}

}

func hand(w http.ResponseWriter, r *http.Request) {
	font(w, r)
	url := r.URL.Path
	h, err := t.GetHand(url, r.Method)
	if err != nil {
		fmt.Printf("警告:该路由不存在 %s", url)
		return
	}
	h(w, r)
	return
}

// 前置处理
func font(w http.ResponseWriter, r *http.Request) {
	// 权限验证
	// 缓存处理
	for _, val := range fontList {
		val(w, r)
	}
}

// AuthorizationDigest 测试http 摘要认证
func AuthorizationDigest(w http.ResponseWriter, r *http.Request) {
	// 获取authorization的头部信息
	auth := r.Header.Get("Authorization")
	fmt.Println("Authorization:", auth)
	fmt.Println("请求方式:", r.Method)
	if auth == "" {
		m := make(map[string]string)
		m["realm"] = "'authrization',"
		nonce := util.StringRand(12)
		m["nonce"] = "'" + nonce + "',"
		m["opaque"] = "'" + util.StringRand(12) + "',"

		fmt.Println("随机字符串:", nonce)

		fmt.Println("WWW-Authenticate:", "Digest "+util.MapToString(m, "="))
		w.Header().Set("WWW-Authenticate",
			"Digest "+util.MapToString(m, "="))

		w.WriteHeader(401)
		w.Write([]byte("你没有权限访问"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("欢迎访问"))
	return

}

var fontList []http.HandlerFunc
var t *tree

func init() {
	t = &tree{make(map[string]*leaf)}
	fontList = make([]http.HandlerFunc, 0)
}
