package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrorNotExsit 记录不存在
	ErrorNotExsit = errors.New("没有找到相应的数据")

	// ErrorInvalid 数据失效错误
	ErrorInvalid = errors.New("缓存的数据已失效")
)

// memory 内存缓存结构
type memory struct {
	// 缓存的数据源
	data []byte

	// 过期时间
	// 从1970.1.1开始的纳秒数(time.Nanosecond)
	expertime int64
}

// 缓存数据是否已过期
// true:未过期
// false: 过期
func (m *memory) isValid() bool {
	return time.Now().Unix()*int64(time.Second) < m.expertime
}

// MemoryCache 内存缓存实例
type MemoryCache struct {
	data map[string]*memory

	maxSize int64
	m       sync.Mutex
}

var _ Cache = new(MemoryCache)

// Get 获取缓存内容
func (mc *MemoryCache) Get(key string) (data []byte, err error) {
	val, ok := mc.data[key]
	if !ok {
		return nil, ErrorNotExsit
	}

	if !val.isValid() {
		return nil, ErrorInvalid
	}
	data = val.data
	return
}

// Put 接受缓存数据
func (mc *MemoryCache) Put(key string, data []byte, t time.Duration) (err error) {
	mc.m.Lock()
	defer mc.m.Unlock()

	val, ok := mc.data[key]
	if ok {
		val.data = data
	} else {
		mc.data[key] = &memory{data: data}
	}
	mc.data[key].expertime = time.Now().Unix()*int64(time.Second) + int64(t)
	return
}

// Delete 删除缓存数据
func (mc *MemoryCache) Delete(key string) error {
	mc.m.Lock()
	defer mc.m.Unlock()

	if _, ok := mc.data[key]; !ok {
		return ErrorNotExsit
	}
	delete(mc.data, key)
	return nil
}

// ClearAll 清空缓存数据
func (mc *MemoryCache) ClearAll() error {
	mc.m.Lock()
	defer mc.m.Unlock()

	mc.data = make(map[string]*memory)
	return nil
}

// IsExist 缓存数据是否存在
func (mc *MemoryCache) IsExist(key string) bool {
	_, ok := mc.data[key]
	return ok
}

// IsFull 是否以达到缓存限制
// 没有找到计算缓存数据大小的方法,暂时允许不限量接收缓存数据
func (mc *MemoryCache) IsFull() bool {
	return len(mc.data) > 20
}

// gc 开启清理无效数据
// 调用时必须启用协程
func (mc *MemoryCache) gc() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			mc.deleteExperted()
		}
	}
}

// deleteExperted 清理无效数据
func (mc *MemoryCache) deleteExperted() {
	for key, val := range mc.data {
		if !val.isValid() {
			mc.m.Lock()
			delete(mc.data, key)
			mc.m.Unlock()
		}
	}
}

// 注册内存缓存到缓存包
func init() {
	m := &MemoryCache{
		data: make(map[string]*memory),
		m:    sync.Mutex{},
	}

	go m.gc()

	RegisteCache(MEMORY, m)
}
