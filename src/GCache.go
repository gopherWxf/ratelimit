package src

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type cacheDate struct {
	key      string
	val      interface{}
	expireAt time.Time
}

func newCacheDate(key string, val interface{}) *cacheDate {
	return &cacheDate{key: key, val: val}
}

type GCache struct {
	elist   *list.List //LRU
	edata   map[string]*list.Element
	lock    sync.Mutex
	maxsize int
}

func NewGCache() *GCache {
	return &GCache{elist: list.New(), edata: make(map[string]*list.Element), maxsize: 100}
}
func (this *GCache) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		this.elist.MoveToFront(v)
		return v.Value.(*cacheDate).val
	}
	return nil
}
func (this *GCache) Set(key string, newVal interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	newCache := newCacheDate(key, newVal)
	if v, ok := this.edata[key]; ok {
		v.Value = newCache
		this.elist.MoveToFront(v)
		return
	} else {
		this.edata[key] = this.elist.PushFront(newCache)
	}
}
func (this *GCache) Print() {
	ele := this.elist.Front()
	if ele == nil {
		return
	}
	for {
		fmt.Println(ele.Value.(*cacheDate).val)
		ele = ele.Next()
		if ele == nil {
			fmt.Println("------")
			break
		}
	}
}
func (this *GCache) RemoveOldest() {
	this.lock.Lock()
	defer this.lock.Unlock()
	back := this.elist.Back()
	if back == nil {
		return
	}
	this.removeItem(back)
}
func (this *GCache) removeItem(ele *list.Element) {
	key := ele.Value.(*cacheDate).key
	delete(this.edata, key)
	this.elist.Remove(ele)
}
