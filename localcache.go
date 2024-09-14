// Package localcache provide a local cache
package localcache

// Cache type describe functionality provided by local cache
type Cache interface {
	// Get value from key
	Get(key string) (value interface{})
	// Set add or replace value of specific key
	Set(key string, value interface{})
	// Fetch will call lambda and lock key if value {nil} ps: even key exists
	Fetch(key string, labmda func() interface{}) (value interface{})
}
