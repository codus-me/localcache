package localcache

import (
	"time"
)

// New will create and return a implementation of Cache
func New() Cache {
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
	cachedData := obj.hashMap[key]
	if cachedData == nil {
		return nil
	}
	outdated := time.Now().After(cachedData.createdAt.Add(ttl))
	if outdated {
		return nil
	}
	return cachedData.data
}

func (obj *cacheImpl) Set(key string, value interface{}) {
	obj.hashMap[key] = &cachedData{
		data:      value,
		createdAt: time.Now(),
	}
}
