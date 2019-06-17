package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
)

var s *Service

type user struct {
	Name string
	Age  int
}

var users []user

func TestMain(m *testing.M) {
	// 每个测试方法运行前都会执行
	users = []user{
		user{"tom", 19},
		user{"ketty", 20},
	}
	s = NewServer()

	go s.RunHTTP(":8081")
	defer s.Shutdown(context.Background())

	code := m.Run()
	os.Exit(code)
}

// 测试基本使用
func TestService(t *testing.T) {
	// 启动服务
	s.GET("/userList", userList)
	url := "http://127.0.0.1:8081/userList"
	err := checkUserList(url)
	if err != nil {
		t.Error(err)
	}

	// 测试自动添加的请求方法验证是否生效
	resp, err := http.Post(url, "", nil)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) != ErrMethodNotAllow {
		t.Error("自动添加的请求方法验证未生效")
	}
}

// 测试分组层级路由
func TestGroup(t *testing.T) {
	g := s.Group("/user")
	g.GET("/list", userList)
	err := checkUserList("http://127.0.0.1:8081/user/list")
	if err != nil {
		t.Error(err)
	}
}

func userList(c *Context) {
	list := make([]user, 0, len(users))
	list = users
	c.JSON(&list)
}

func checkUserList(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	list := make([]user, 0, len(users))
	err = json.Unmarshal(body, &list)
	if err != nil {
		return err
	}
	if len(list) != len(users) {
		return errors.New("GET请求获取的数据失败")
	}
	return nil
}

// 测试 gratefully Shutdown
func TestShutdown(t *testing.T) {
	ch := make(chan struct{})
	s.GET("/timeout", func(c *Context) {
		ch <- struct{}{}
		time.Sleep(5 * time.Second)
		c.Writer.Write([]byte("休息五秒"))
	})

	f := func() {
		resp, err := http.Get("http://127.0.0.1:8081/timeout")
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		t.Log(string(body))
	}
	go f()
	<-ch
	t.Log("grateful shutdown")
	err := s.Shutdown(context.Background())
	if err != nil {
		t.Error(err)
	}
}
