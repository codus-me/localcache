// Package localcache provide a local cache
package localcache

import "time"

// Cache type describe functionality provided by local cache
type Cache interface {
	// Get value from key
	Get(key string) (value interface{})
	// Set add or replace value of specific key
	Set(key string, value interface{})
}

// TTL time to live 30 seconds
const TTL time.Duration = 30 * time.Second

// New will create and return a implementation of Cache
func New() Cache {
	return newCacheImpl()
}
