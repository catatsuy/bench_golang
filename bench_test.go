package main

import (
	"sync"
	"testing"
)

type cache struct {
	sync.RWMutex
	items map[string]interface{}
}

func New() *cache {
	m := make(map[string]interface{})
	c := &cache{
		items: m,
	}
	return c
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.RLock()
	v, found := c.items[key]
	c.RUnlock()
	return v, found
}

func (c *cache) Set(key string, value interface{}) {
	c.Lock()
	c.items[key] = value
	c.Unlock()
}

func BenchmarkCacheGetHit(b *testing.B) {
	b.StopTimer()
	tc := New()
	tc.Set("foo", "bar")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foo")
	}
}

func BenchmarkCacheGetNoHit(b *testing.B) {
	b.StopTimer()
	tc := New()
	tc.Set("foo", "bar")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tc.Get("foo1")
	}
}

type TestStruct struct {
	Num int
}

type cacheStruct struct {
	sync.RWMutex
	items map[string]TestStruct
}

func NewStruct() *cacheStruct {
	m := make(map[string]TestStruct)
	c := &cacheStruct{
		items: m,
	}
	return c
}

func (c *cacheStruct) Get(key string) (TestStruct, bool) {
	c.RLock()
	v, found := c.items[key]
	c.RUnlock()
	return v, found
}

func (c *cacheStruct) Set(key string, value *TestStruct) {
	c.Lock()
	c.items[key] = *value
	c.Unlock()
}

func BenchmarkCacheGetHitStruct(b *testing.B) {
	b.StopTimer()
	tc := New()
	tc.Set("foo", &TestStruct{Num: 1})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t, _ := tc.Get("foo")
		_ = t.(*TestStruct)
	}
}

func BenchmarkCacheGetHitTestStruct(b *testing.B) {
	b.StopTimer()
	tc := NewStruct()
	tc.Set("foo", &TestStruct{Num: 1})
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		t, _ := tc.Get("foo")
		_ = t
	}
}
