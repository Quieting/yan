package yan

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// 通过寻找配置文件自动注册路由,路由规则:
// URI最后一级对应方法,倒数第二级对应结构体,多余层级路由全是包名
// 因此URI至少是两级URI

// Hand http业务处理函数
type Hand func(w http.ResponseWriter, r *http.Request)

// 错误变量
var (
	// 没有找到与url匹配的处理方法
	ErrMismatch = errors.New("URI is not match")
)

type tree struct {
	leafs map[string]*leaf
}

type leaf struct {
	leafs    map[string]*leaf
	fullPath string
	path     string
	methods  []string // len(l.Methods) == 0 表示全匹配
	h        Hand
}

// RegisterRouter 注册路由
// path:完整的URI(不包含URI参数),前后没有 "/" 的格式
func (t *tree) RegisterRouter(path string, f Hand, methods ...string) {
	// 以 "/" 分割成数组
	paths := strings.Split(path, "/")

	if _, ok := t.leafs[paths[0]]; !ok {
		t.leafs[paths[0]] = newLeaf()
	}
	t.leafs[paths[0]].addLeaf(paths, f, methods...)
}

// addLeaf 添加路由枝叶
// 方法解析:使用递归添加路径,递归一次添加一个路由节点
func (l *leaf) addLeaf(paths []string, f Hand, methods ...string) {
	if len(paths) == 1 {
		if l.leafs == nil {
			l.leafs = make(map[string]*leaf)
		}

		if _, ok := l.leafs[paths[0]]; ok {
			// 覆盖原映射的方法
			l.leafs[paths[0]].h = f
			l.methods = methods
			fmt.Println("Waring:重复添加路由已覆盖")
		}

		// 把http请求方法转换成大写
		if len(methods) > 0 {
			for i := 0; i < len(methods); i++ {
				methods[i] = strings.ToUpper(methods[i])
			}
		}

		l.leafs[paths[0]] = &leaf{
			path:     paths[0],
			fullPath: l.fullPath + "/" + paths[0],
			methods:  methods,
			h:        f,
		}

		fmt.Println("l.leafs[paths[0]]:", l.leafs[paths[0]])
		return
	}

	if _, ok := l.leafs[paths[0]]; !ok {
		l.leafs[paths[0]] = newLeaf()
		l.leafs[paths[0]].fullPath = l.fullPath + "/" + paths[0]
	}

	l.leafs[paths[0]].addLeaf(paths[1:], f, methods...)
}

// getHand 通过路由获取处理方法
func (t *tree) GetHand(path, method string) (f Hand, err error) {
	// 以 "/" 分割成数组
	paths := make([]string, 0, 5)
	if strings.HasPrefix(path, "/") {

	}
	paths = strings.Split(path, "/")
	if _, ok := t.leafs[paths[0]]; !ok {
		return nil, ErrMismatch
	}
	return t.leafs[paths[0]].getLeaf(paths, method)
}

// getLeaf 遍历路由枝叶
// 方法解析:使用递归解析路径,最后判断请求方法是否匹配
func (l *leaf) getLeaf(paths []string, method string) (f Hand, err error) {
	// 检查最后一级路径是否匹配
	if len(paths) == 1 {
		if _, ok := l.leafs[paths[0]]; !ok {
			return nil, ErrMismatch
		} else if len(l.leafs[paths[0]].methods) > 0 {
			// 检查请求方式是否匹配,len(l.Methods) == 0 表示全匹配
			methodIsMatch := false
			for _, val := range l.leafs[paths[0]].methods {
				if strings.ToUpper(method) == val {
					methodIsMatch = true
				}
			}
			if !methodIsMatch {
				return nil, ErrMismatch
			}
		}

		// 获取处理方法
		if l.leafs[paths[0]].path != "" && l.leafs[paths[0]].path == paths[0] {
			fmt.Println("fullPath:", l.leafs[paths[0]].fullPath)
			return l.leafs[paths[0]].h, nil
		}
		return nil, ErrMismatch
	}

	if _, ok := l.leafs[paths[0]]; !ok {
		return nil, ErrMismatch
	}

	return l.leafs[paths[0]].getLeaf(paths[1:], method)
}

// 得到一个新的枝叶
func newLeaf() *leaf {
	return &leaf{
		leafs:   make(map[string]*leaf),
		methods: make([]string, 0, 10),
	}
}
