package src

import (
	"container/list"
	"fmt"
	"sync"
)

type GCache struct {
	elist *list.List //LRU
	edata map[string]*list.Element
	lock  sync.Mutex
}

func NewGCache() *GCache {
	return &GCache{elist: list.New(), edata: make(map[string]*list.Element)}
}
func (this *GCache) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		this.elist.MoveToFront(v)
		return v.Value
	}
	return nil
}
func (this *GCache) Set(key string, newVal interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if v, ok := this.edata[key]; ok {
		v.Value = newVal
		this.elist.MoveToFront(v)
		return
	} else {
		this.edata[key] = this.elist.PushFront(newVal)
	}
}
func (this *GCache) Print() {
	ele := this.elist.Front()
	if ele == nil {
		return
	}
	for {
		fmt.Println(ele.Value)
		ele = ele.Next()
		if ele == nil {
			fmt.Println("------")
			break
		}
	}
}
