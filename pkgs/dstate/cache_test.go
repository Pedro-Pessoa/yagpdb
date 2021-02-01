package dstate

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newCache() *Cache {
	return NewCache(new(int64), new(int64))
}

func TestCacheSetGet(t *testing.T) {
	c := newCache()

	key := "123"
	c.Set(key, "hey")

	v := c.Get(key)
	if v == nil {
		t.Error("value did not set")
	}

	if v != "hey" {
		t.Error("value is not 'hey': ", v)
	}
}

func TestCacheEviction(t *testing.T) {
	c := newCache()

	key := "123"
	c.Set(key, "hey")

	n := c.EvictOldKeys(time.Now().Add(time.Hour))
	if n != 1 {
		t.Error("did not evict any keys")
		return
	}

	v := c.Get(key)
	if v != nil {
		t.Error("value is not nil after eviction")
	}
}

func TestCacheConcurrentSingleFetch(t *testing.T) {
	c := newCache()

	key := "123"

	didFetch := new(int64)

	var wg sync.WaitGroup

	fetchFunc := func() (interface{}, error) {
		t.Log("In fetch")

		if !atomic.CompareAndSwapInt64(didFetch, 0, 1) {
			t.Error("Concurrent fetch calls!")
			return 2, nil
		}

		defer wg.Done()

		time.Sleep(time.Second * 2)
		return 1, nil
	}

	wg.Add(1)

	go func() { _, _ = c.Fetch(key, fetchFunc) }()
	v, err := c.Fetch(key, fetchFunc)
	if err != nil {
		t.Log("Error fetching: ", err)
	}

	if v != 1 {
		t.Log("Incorrect value")
	}

	wg.Wait()
}

func TestCachePanicRecovery(t *testing.T) {
	c := newCache()

	key := "123"

	didFetch := new(int64)

	fetchFunc := func() (interface{}, error) {
		if !atomic.CompareAndSwapInt64(didFetch, 0, 1) {
			return 2, nil
		}

		t.Log("Panicing")
		panic("wew")
	}

	_, _ = fetchRecover(c, fetchFunc, key)

	v, err := c.Fetch(key, fetchFunc)
	if err != nil {
		t.Log("Error fetching: ", err)
	}

	if v != 2 {
		t.Log("Incorrect value")
	}
}

func TestCachePanicRecoveryWaiting(t *testing.T) {
	c := newCache()

	key := "123"

	didFetch := new(int64)

	fetchFunc := func() (interface{}, error) {
		if !atomic.CompareAndSwapInt64(didFetch, 0, 1) {
			return 2, nil
		}

		time.Sleep(time.Second * 3)
		t.Log("Panicing")
		panic("wew")
	}

	go func() { _, _ = fetchRecover(c, fetchFunc, key) }()
	time.Sleep(time.Second)
	v, err := c.Fetch(key, fetchFunc)
	if err != nil {
		t.Log("Error fetching: ", err)
	}

	if v != 2 {
		t.Log("Incorrect value")
	}
}

func fetchRecover(c *Cache, fetchFunc CacheFetchFunc, key string) (interface{}, error) {
	defer func() {
		_ = recover()
	}()

	return c.Fetch(key, fetchFunc)
}

type AT int
type BT int

func TestCacheDelAllType(t *testing.T) {
	a1 := AT(1)
	a2 := AT(2)
	b1 := BT(1)

	c := newCache()

	c.Set(a1, 1)
	c.Set(a2, 1)
	c.Set(b1, 1)
	c.Set(1, 1)

	c.DelAllKeysType(a1)

	if c.Get(a1) != nil || c.Get(a2) != nil {
		t.Error("Should have been removed from cache")
	}

	if c.Get(b1) == nil {
		t.Error("Should still exist in cache")
	}

	if c.Get(1) == nil {
		t.Error("Should still exist in cache")
	}
}
