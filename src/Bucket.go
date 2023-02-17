package src

import (
	"sync"
	"time"
)

type Bucket struct {
	cap    int
	tokens int
	lock   sync.Mutex
	rate   int
}

func NewBucket(cap int, rate int) *Bucket {
	if cap < 0 || rate < 0 {
		panic("err cap")
	}
	bucket := &Bucket{cap: cap, tokens: cap, rate: rate}
	bucket.start()
	return bucket
}
func (this *Bucket) start() {
	go func() {
		for {
			time.Sleep(time.Second * 1)
			this.addToken()
		}
	}()
}
func (this *Bucket) addToken() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.tokens += this.rate
	if this.tokens > this.rate {
		this.tokens = this.rate
	}
}
func (this *Bucket) IsAccept() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.tokens > 0 {
		this.tokens--
		return true
	}
	return false
}
