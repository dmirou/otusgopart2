package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		// Write me
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove if task with asterisk completed

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1000000))))
		}
	}()

	wg.Wait()
}

// TestSet check simple setting elements to cache.
func TestSet(t *testing.T) {
	c := NewCache(3)

	updated := c.Set("key1", 2)
	if updated {
		t.Errorf("expected updated: %v, got: %v", false, updated)
	}

	updated = c.Set("key2", "hello")
	if updated {
		t.Errorf("expected updated: %v, got: %v", false, updated)
	}

	v, ok := c.Get("key1")
	if !ok {
		t.Errorf("expected ok: %v, got: %v", true, ok)
	}

	vInt, ok := v.(int)
	if !ok {
		t.Errorf("can not assert value type as int for value: %v", v)
	}

	if vInt != 2 {
		t.Errorf("expected vInt: %v, got: %v", 2, vInt)
	}

	v, ok = c.Get("key2")
	if !ok {
		t.Errorf("expected ok: %v, got: %v", true, ok)
	}

	vStr, ok := v.(string)
	if !ok {
		t.Errorf("can not assert value type as string for value: %v", v)
	}

	if vStr != "hello" {
		t.Errorf("expected vStr: %v, got: %v", "hello", vStr)
	}
}

// TestSetPop checks that last element will be
// deleted from the cache if cache size is exceeded.
func TestSetPop(t *testing.T) {
	c := NewCache(2)

	c.Set("k1", 4)
	c.Set("k2", 5)
	c.Set("k3", 9)

	v, ok := c.Get("k1")
	if ok {
		t.Errorf("expected ok: %v, got: %v", false, ok)
	}

	if v != nil {
		t.Errorf("expected v: %v, got: %v", nil, v)
	}

	v, ok = c.Get("k2")
	if !ok {
		t.Errorf("expected ok: %v, got: %v", true, ok)
	}

	if v != 5 {
		t.Errorf("expected v: %v, got: %v", 5, v)
	}

	v, ok = c.Get("k3")
	if !ok {
		t.Errorf("expected ok: %v, got: %v", true, ok)
	}

	if v != 9 {
		t.Errorf("expected v: %v, got: %v", 9, v)
	}
}

// TestSetPopLeastUsed checks that the least used
// item will be deleted from the queue to set
// new item if cache size is exceeded.
func TestSetPopLeastUsed(t *testing.T) {
	c := NewCache(3)

	c.Set("k1", 4)
	c.Set("k2", 5)
	c.Set("k3", 9)

	_, ok := c.Get("k3")
	if !ok {
		t.Errorf("expected ok: %v, got: %v", true, ok)
	}

	updated := c.Set("k2", 13)
	if !updated {
		t.Errorf("expected updated: %v, got: %v", true, updated)
	}

	updated = c.Set("k3", 2)
	if !updated {
		t.Errorf("expected updated: %v, got: %v", true, updated)
	}

	updated = c.Set("k4", 90)
	if updated {
		t.Errorf("expected updated: %v, got: %v", false, updated)
	}

	v, ok := c.Get("k1")
	if ok {
		t.Errorf("expected ok: %v, got: %v", false, ok)
	}

	if v != nil {
		t.Errorf("expected v: %v, got: %v", nil, v)
	}
}
