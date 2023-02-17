package src

import "testing"

func TestNewGCache(t *testing.T) {
	cache := NewGCache()
	//LRU
	cache.Set("name", "wxf")
	cache.Set("age", 20)
	cache.Set("sex", "男")
	cache.Print()

	cache.Get("name")
	cache.Set("age", 19)
	cache.Print()

	cache.RemoveOldest()
	cache.Print()
}
