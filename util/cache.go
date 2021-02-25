package util

import (
	"github.com/smallsixFight/hyakkei_blog/task"
	"sync"
	"time"
)

type customCache struct {
	m      sync.Map
	record sync.Map
}

var Cache *customCache

func (c *customCache) Set(key, val interface{}, duration time.Duration) {
	if c.IsExist(key) {
		c.Del(key)
	}
	c.m.Store(key, val)
	// 不为 0 则设置定时删除任务
	if duration != 0 {
		c.record.Store(key, task.Scheduler.AddTask(func() {
			c.Del(key)
		}, duration))
	}
}

func (c *customCache) Del(key interface{}) {
	v, ok := c.record.Load(key)
	c.m.Delete(key)
	if ok {
		task.Scheduler.RemoveTask(v.(int64))
	}
}

func (c *customCache) Get(key interface{}) interface{} {
	val, _ := c.m.Load(key)
	return val
}

func (c *customCache) GetBool(key interface{}) bool {
	val, _ := c.m.Load(key)
	if v, ok := val.(bool); ok {
		return v
	}
	return false
}

func (c *customCache) GetBytes(key interface{}) []byte {
	val, _ := c.m.Load(key)
	if v, ok := val.([]byte); ok {
		return v
	}
	return nil
}

func (c *customCache) GetString(key interface{}) string {
	val, _ := c.m.Load(key)
	if v, ok := val.(string); ok {
		return v
	}
	return ""
}

//func (c *customCache) GetInt(key interface{}) int {
//	val, _ := c.m.Load(key)
//	if v, ok := val.(int); ok {
//		return v
//	}
//	return 0
//}

func (c *customCache) GetInt32(key interface{}) int32 {
	val, _ := c.m.Load(key)
	if v, ok := val.(int32); ok {
		return v
	}
	return 0
}

func (c *customCache) GetInt64(key interface{}) int64 {
	val, _ := c.m.Load(key)
	if v, ok := val.(int64); ok {
		return v
	}
	return 0
}

func (c *customCache) GetAndDel(key interface{}) interface{} {
	val, ok := c.m.Load(key)
	if ok {
		c.Del(key)
	}
	return val
}

func (c *customCache) IsExist(key interface{}) bool {
	_, ok := c.m.Load(key)
	return ok
}

func init() {
	Cache = &customCache{m: sync.Map{}, record: sync.Map{}}
}
