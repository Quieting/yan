/*
 *
 *
 *
 *
 */
package cache

import (
	"errors"
	"time"
)

// Cache 每个缓存实例必须实现的接口
type Cache interface {
	// 获取缓存内容,如果缓存不存在返回 ErrorIsNotExist 错误
	Get(key string) ([]byte, error)

	// 缓存内容,支持并发
	// 缓存key已存在时,更换缓存内容并重新计算缓存时间
	// key:数据唯一标识符
	// data:缓存数据
	// t:数据缓存时间
	Put(key string, data []byte, t time.Duration) error

	// 删除缓存,支持并发
	Delete(key string) error

	// 清除所有缓存
	ClearAll() error

	// 检查缓存内容是否存在,存在返回true,不存在返回false
	IsExist(key string) bool

	// 检查缓存是否已存满,存满返回true,没存满返回false
	IsFull() bool
}

// 缓存实例适配器使用前需要先注册
var cacheAdapt map[medium]Cache

type strategy int

// 缓存策略
const (
	// FIFO 先进先出策略
	FIFO strategy = 1 << iota

	// LFU 最少使用
	LFU

	// LRU 最久未被使用
	LRU
)

type medium string

// 缓存介质
const (
	// MEMORY 内存缓存
	MEMORY medium = "memory"
)

// 使用到的错误
var (
	ErrorIsExist         = errors.New("注册的缓存实例key已存在")
	ErrorIsNotExist      = errors.New("未找到缓存实例")
	ErrorInvalidStrategy = errors.New("无效的strategy值")
)

// RegisteCache 注册自定义缓存实例
func RegisteCache(key medium, c Cache) error {
	if _, ok := cacheAdapt[key]; ok {
		return ErrorIsExist
	}
	cacheAdapt[key] = c
	return nil
}

// NewCacherWithStrategy 获取一个缓存实例带有缓存策略
// 若果未找到缓存实例则返回错误
// key:缓存的媒介(支持的值:memory)
// strategy:缓存的策略
func NewCacherWithStrategy(key medium, stra strategy) (Cache, error) {
	c, ok := cacheAdapt[key]
	if !ok {
		return nil, ErrorIsNotExist
	}
	switch stra {
	case FIFO:
		return NewFIFOStrategy(c), nil
	}

	return nil, ErrorInvalidStrategy
}

// Cacher 使用的缓存实例
type Cacher struct {
	Cache
	Strategy int
}

func init() {
	cacheAdapt = make(map[medium]Cache)
}
