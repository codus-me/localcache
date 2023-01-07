package localcache

import (
	"sync"
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
	mux     sync.RWMutex
}

type cachedData struct {
	data      interface{}
	createdAt time.Time
}

func (obj *cacheImpl) Get(key string) interface{} {
	obj.mux.RLock()
	defer obj.mux.RUnlock()
	cachedData := obj.hashMap[key]
	if cachedData == nil {
		return nil
	} else if ok := time.Now().Before(cachedData.createdAt.Add(ttl)); !ok {
		return nil
	}
	return cachedData.data
}

func (obj *cacheImpl) Set(key string, value interface{}) {
	obj.mux.Lock()
	defer obj.mux.Unlock()
	obj.hashMap[key] = &cachedData{
		data:      value,
		createdAt: time.Now(),
	}
}
