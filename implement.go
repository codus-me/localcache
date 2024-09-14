package localcache

import (
	"log/slog"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// New will create and return a implementation of Cache
func New() Cache {
	return &cacheImpl{
		hashMap: make(map[string]*cachedData),
		lockMap: make(map[string]*sync.Mutex),
	}
}

const ttl time.Duration = 30 * time.Second

type cacheImpl struct {
	hashMap map[string]*cachedData
	lockMap map[string]*sync.Mutex // lockmap 用於 fetch
	engine  singleflight.Group
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

// Do 相較於 fetch，使用了 golang singleflight 來避免重複呼叫
func (obj *cacheImpl) Do(key string, lambda func() interface{}) (value interface{}, err error) {
	value = obj.Get(key)
	if value != nil {
		return
	}
	row, err, _ := obj.engine.Do(key, func() (interface{}, error) {
		data := lambda()
		return data, nil
	})

	if err != nil {
		slog.Error("singleflight error", "err", err)
		return nil, err
	}
	return row, nil
}

// Fetch 會避免 miss key，重複呼叫，重複執行。自行實作. 而Do 使用 golang singleflight
func (obj *cacheImpl) Fetch(key string, lambda func() interface{}) (value interface{}) {
	fetchLock := obj.getFetchLock(key)

	fetchLock.Lock()
	defer fetchLock.Unlock()

	value = obj.Get(key)
	if value != nil {
		return
	}
	value = lambda()
	obj.Set(key, value)
	return
}

func (obj *cacheImpl) getFetchLock(key string) (fetchLock *sync.Mutex) {
	obj.mux.Lock()
	fetchLock = obj.lockMap[key]
	if fetchLock == nil {
		fetchLock = &sync.Mutex{}
		obj.lockMap[key] = fetchLock
	}
	obj.mux.Unlock()
	return
}
