package xcache

import (
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

// Cache 生命缓存实例
var Cache *cache.Cache

// NewCache 创建缓存实例
func NewCache(defaultExpiration, cleanupTime time.Duration) *cache.Cache {
	var once sync.Once
	once.Do(func() {
		Cache = cache.New(defaultExpiration, cleanupTime)
	})
	return Cache
}

// NewDefaultCache 创建默认的缓存实例
func NewDefaultCache() *cache.Cache {
	return NewCache(time.Minute*5, time.Minute*10)
}


