// 缓存先进先出策略实现
package cache

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"
)

// FIFOStrategy 先入先出缓存策略
type FIFOStrategy struct {
	c        Cache
	keyArray []string

	// 值无效的下标
	invalidIndex int

	m sync.Mutex
}

// Get 获取缓存内容
func (f *FIFOStrategy) Get(key string) (data []byte, err error) {
	return f.c.Get(key)
}

// Put 缓存内容
func (f *FIFOStrategy) Put(key string, data []byte, t time.Duration) (err error) {
	f.m.Lock()
	defer f.m.Unlock()

	if f.IsFull() {
		f.Delete(f.keyArray[f.invalidIndex])
		err = f.c.Put(key, data, t)
		if err != nil {
			return
		}
		f.invalidIndex++

		if f.invalidIndex >= len(f.keyArray) {
			f.invalidIndex = 0
		}
		return
	}

	err = f.c.Put(key, data, t)
	if err != nil {
		return
	}

	f.keyArray = append(f.keyArray, key)
	return
}

// Delete 删除缓存数据
func (f *FIFOStrategy) Delete(key string) error {
	return f.c.Delete(key)
}

// ClearAll 清空缓存数据
func (f *FIFOStrategy) ClearAll() error {
	err := f.c.ClearAll()
	if err != nil {
		return err
	}
	f.keyArray = make([]string, 0, 10)
	f.invalidIndex = 0
	return nil
}

// IsExist 检查缓存对象是否存在
func (f *FIFOStrategy) IsExist(key string) bool {
	return f.c.IsExist(key)
}

// IsFull 检查缓存是否已满
func (f *FIFOStrategy) IsFull() bool {
	return f.c.IsFull()
}

// NewFIFOStrategy 获取一个先入先出缓存策略对象
func NewFIFOStrategy(c Cache) *FIFOStrategy {
	return &FIFOStrategy{
		c:        c,
		keyArray: make([]string, 0, 10),
	}
}

// Encode 使用内存缓存时获取缓存对象的内存大小有两种方式
// 1.通过gob对数据结构重新编码,可获取具体的大小但是效率较低
// 2.调用unsafe包的SizeOf方法,需要明确数据类型,并且无法获取指针对象的实际大小
// 3.使用反射
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
