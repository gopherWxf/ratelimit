package src

import (
	"sync"
	"time"
)

type Bucket struct {
	cap      int64
	tokens   int64
	lock     sync.Mutex
	rate     int64 //每秒加入token数量
	lastTime int64
}

func NewBucket(cap int64, rate int64) *Bucket {
	if cap < 0 || rate < 0 {
		panic("err cap")
	}
	bucket := &Bucket{cap: cap, tokens: cap, rate: rate}
	return bucket
}

func (this *Bucket) IsAccept() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	now := time.Now().Unix()
	this.tokens += (now - this.lastTime) * this.rate
	if this.tokens > this.cap {
		this.tokens = this.cap
	}
	this.lastTime = now

	if this.tokens > 0 {
		this.tokens--
		return true
	}
	return false
}
