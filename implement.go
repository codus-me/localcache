package localcache

import (
	"time"
)

type cacheImpl struct {
	data      map[string]interface{}
	createdAt map[string]*time.Time
}

func newCacheImpl() *cacheImpl {
	return &cacheImpl{
		data:      make(map[string]interface{}),
		createdAt: make(map[string]*time.Time),
	}
}

func (obj *cacheImpl) Get(key string) interface{} {
	if obj.createdAt[key] != nil && time.Now().After(obj.createdAt[key].Add(TTL)) {
		return nil
	}
	return obj.data[key]
}

func (obj *cacheImpl) Set(key string, value interface{}) {
	now := time.Now()
	obj.createdAt[key] = &now
	obj.data[key] = value
}
