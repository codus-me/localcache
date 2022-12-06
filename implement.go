package localcache

import (
	"time"
)

// New will create and return a implementation of Cache
func New() *cacheImpl {
	return &cacheImpl{
		hashMap: make(map[string]*cachedData),
	}
}

const ttl time.Duration = 30 * time.Second

type cacheImpl struct {
	hashMap map[string]*cachedData
}

type cachedData struct {
	data      interface{}
	createdAt time.Time
}

func (obj *cacheImpl) Get(key string) interface{} {
	if obj.hashMap[key] == nil {
		return nil
	}
	outdated := time.Now().After(obj.hashMap[key].createdAt.Add(ttl))
	if outdated {
		return nil
	}
	return obj.hashMap[key].data
}

func (obj *cacheImpl) Set(key string, value interface{}) {
	obj.hashMap[key] = &cachedData{
		data:      value,
		createdAt: time.Now(),
	}
}
