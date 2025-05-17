package cache

import "github.com/jellydator/ttlcache/v3"

var cache *ttlcache.Cache[string, any]

func Active() *ttlcache.Cache[string, any] {
	if cache == nil { // coverage-ignore
		panic("cache is not initialized")
	}

	return cache
}

func Init() {
	if cache != nil {
		cache.Stop()
	}

	cache = ttlcache.New[string, any]()

	go cache.Start()
}

func Close() {
	cache.Stop()
}
