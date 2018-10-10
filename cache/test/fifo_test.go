package test

import (
	"fmt"
	"testing"
	"time"

	"http/cache"
)

var f cache.Cache

func init() {
	r, err := cache.NewCacherWithStrategy("memory", cache.FIFO)
	if err == nil {
		f = r
	} else {
		panic(err)
	}
}

// 测试先进先出策略是否生效
func Test_FIFO(t *testing.T) {
	// 缓存数据
	for i := 1; i < 30; i++ {
		err := f.Put(fmt.Sprintf("%d", i), []byte(fmt.Sprintf("我是第%d个缓存数据", i)), time.Second*1)
		if err != nil {
			t.Error(err)
		}

		if i > 20 {
			da, err := f.Get(fmt.Sprintf("%d", 1-19))
			if err != cache.ErrorNotExsit {
				t.Log("数据未清除:", da)
			}
			da, err = f.Get(fmt.Sprintf("%d", i))
			if string(da) != fmt.Sprintf("我是第%d个缓存数据", i) {
				t.Log("数据异常")
			}
		}
	}
}

// 测试策略效率
func Benchmark_FIFOSpeed(t *testing.B) {
	for i := 0; i < t.N; i++ {
		err := f.Put(fmt.Sprintf("%d", i), []byte(fmt.Sprintf("我是第%d个缓存数据", i)), time.Second*1)
		if err != nil {
			t.Error(err)
		}
	}
}
