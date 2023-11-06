package cache

import "github.com/jellydator/ttlcache/v3"

func Init() *ttlcache.Cache[string, string] {
	cache := ttlcache.New[string, string]()

	go cache.Start() // starts automatic expired item deletion

	return cache
}
