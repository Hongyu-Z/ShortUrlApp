package cache

import "github.com/jellydator/ttlcache/v3"

func Init() *ttlcache.Cache[string, string] {
	capacity := 10000
	cache := ttlcache.New[string, string](ttlcache.WithCapacity[string, string](uint64(capacity)))

	go cache.Start() // starts automatic expired item deletion

	return cache
}
