package util

import (
	"testing"
	"time"
)

func TestGetCache(t *testing.T) {
	Cache.Set(1, 2, time.Second*3)
	t.Log(Cache.Get(1))
	time.Sleep(time.Second)
	t.Log(Cache.Get(1))
	time.Sleep(time.Second * 2)
	t.Log(Cache.Get(1))
}

func TestGetCache2(t *testing.T) {
	Cache.Set(1, 2, time.Second*3)
	Cache.Set(1, 3, time.Second)
	t.Log(Cache.Get(1))
	time.Sleep(time.Second)
	t.Log(Cache.Get(1))
}
