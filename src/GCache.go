package src

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type cacheData struct {
	key      string
	val      interface{}
	expireAt time.Time
}

func newCacheDate(key string, val interface{}, expireAt time.Time) *cacheData {
	return &cacheData{key: key, val: val, expireAt: expireAt}
}
func (this *cacheData) IsExpire() bool {
	if time.Now().After(this.expireAt) {
		return true
	}
	return false
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
	cache.clear()
	return cache
}
func (this *GCache) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		//过期了
		if v.Value.(*cacheData).IsExpire() {
			return nil
		}
		this.elist.MoveToFront(v)
		return v.Value.(*cacheData).val
	}
	return nil
}

// 不过期的时间
const NotExpireTTL = time.Hour * 24 * 356

func (this *GCache) Set(key string, newVal interface{}, TTL time.Duration) {
	this.lock.Lock()
	defer this.lock.Unlock()
	var setExpire time.Time
	if TTL == 0 {
		setExpire = time.Now().Add(NotExpireTTL)
	} else {
		setExpire = time.Now().Add(TTL)
	}
	newCache := newCacheDate(key, newVal, setExpire)
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
	for e := this.elist.Front(); e != nil; e = e.Next() {
		fmt.Printf("[key:%v val:%v] ",
			e.Value.(*cacheData).key,
			e.Value.(*cacheData).val)
	}
	fmt.Printf("\n")
}
func (this *GCache) removeOldest() {
	back := this.elist.Back()
	if back == nil {
		return
	}
	this.removeItem(back)
}
func (this *GCache) removeItem(ele *list.Element) {
	key := ele.Value.(*cacheData).key
	delete(this.edata, key)
	this.elist.Remove(ele)
}
func (this *GCache) removeExpired() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, v := range this.edata {
		temp := v.Value.(*cacheData)
		if temp.IsExpire() {
			this.removeItem(v)
		}
	}
}
func (this *GCache) clear() {
	go func() {
		for {
			this.removeExpired()
			time.Sleep(time.Second * 1)
		}
	}()
}
func (this *GCache) len() int {
	return len(this.edata)
}
