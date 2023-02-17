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

type GCacheOption func(g *GCache)
type GCacheOptions []GCacheOption

func (this GCacheOptions) Apply(g *GCache) {
	for _, fn := range this {
		fn(g)
	}
}
func WithMaxSize(size int) GCacheOption {
	return func(g *GCache) {
		if size > 0 {
			g.maxsize = size
		}
	}
}

type GCache struct {
	elist   *list.List //LRU
	edata   map[string]*list.Element
	lock    sync.Mutex
	maxsize int
}

func NewGCache(opt ...GCacheOption) *GCache {
	cache := &GCache{elist: list.New(), edata: make(map[string]*list.Element), maxsize: 0}
	GCacheOptions(opt).Apply(cache)
	return cache
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
		//判断长度是否溢出,溢出则淘汰末位缓存
		if this.maxsize > 0 && len(this.edata) > this.maxsize {
			this.removeOldest()
		}
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
func (this *GCache) removeOldest() {
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
