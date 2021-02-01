package dstate

import (
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	store     map[interface{}]*bucket
	cacheHits *int64
	cacheMiss *int64

	cond *sync.Cond
}

// NewCache creates a new cache
func NewCache(hit, miss *int64) *Cache {
	return &Cache{
		store:     make(map[interface{}]*bucket),
		cacheHits: hit,
		cacheMiss: miss,
		cond:      sync.NewCond(&sync.Mutex{}),
	}
}

type bucket struct {
	Value    interface{}
	Created  time.Time
	Fetching bool
}

// Get retrieves an item from the cache, or nil if non-existing
func (c *Cache) Get(key interface{}) interface{} {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	return c.get(key)
}

func (c *Cache) get(key interface{}) interface{} {
	for {
		bucket := c.store[key]
		if bucket == nil {
			return nil
		}

		if bucket.Fetching {
			c.cond.Wait()
			continue
		}

		return bucket.Value
	}

}

// Set stores an item in the cache, this does mutate and needs to be used with a full lock
func (c *Cache) Set(key, value interface{}) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	c.store[key] = &bucket{
		Created: time.Now(),
		Value:   value,
	}
}

// Del deletes an item from the cache, this mutates state and needs to be used with a full lock, if its used from multiple goroutines
func (c *Cache) Del(key interface{}) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	delete(c.store, key)
}

// DelAllKeysType deletes all keys that are of a specific type
func (c *Cache) DelAllKeysType(t interface{}) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	typ := reflect.TypeOf(t)

	for k := range c.store {
		vt := reflect.TypeOf(k)
		if typ == vt {
			delete(c.store, k)
		}
	}
}

type CacheFetchFunc func() (value interface{}, err error)

// Fetch either retrieves an item directly from the cache if available, or calls the passed fetchFunc to then set it in cache and return it
func (c *Cache) Fetch(key interface{}, fetchFunc CacheFetchFunc) (interface{}, error) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	for {
		bucket := c.store[key]
		if bucket == nil {
			// need to fetch from db
			break
		}

		if bucket.Fetching {
			// wait until were done fetching
			c.cond.Wait()
			continue
		} else {
			// cache hit!
			atomic.AddInt64(c.cacheHits, 1)
			return bucket.Value, nil
		}
	}

	// bucket did not exist, create it and fetch from underlying db

	// wake all wait() calls after we leave the function
	defer c.cond.Broadcast()

	clearBucket := true
	defer func() {
		// clear bucket if we panic inside fetch so that we don't get a deadlock
		if clearBucket {
			delete(c.store, key)
		}
	}()

	atomic.AddInt64(c.cacheMiss, 1)

	// mark as fetching to avoid simoultanous fetch calls
	bucket := &bucket{
		Created:  time.Now(),
		Fetching: true,
	}
	c.store[key] = bucket

	v, err := c.fetchUnderlying(key, fetchFunc)
	if err != nil {
		return nil, err
	}

	bucket.Value = v
	bucket.Fetching = false
	clearBucket = false

	return v, nil
}

func (c *Cache) fetchUnderlying(key interface{}, fetchFunc CacheFetchFunc) (interface{}, error) {
	c.cond.L.Unlock()
	defer c.cond.L.Lock()

	return fetchFunc()
}

// EvictOldKeys evicts old keys created before treshold
func (c *Cache) EvictOldKeys(treshold time.Time) (evicted int) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	for k, v := range c.store {
		if v.Created.Before(treshold) {
			delete(c.store, k)
			evicted++
		}
	}

	return
}
