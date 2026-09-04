package cart_cache

import "github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"

type cartMencache struct {
	CartQueryCache
	CartCommandCache
}

type CartMencache interface {
	CartQueryCache
	CartCommandCache
}

func NewCartMencache(cacheStore *cache.CacheStore) CartMencache {
	return &cartMencache{
		CartQueryCache:  NewCartQueryCache(cacheStore),
		CartCommandCache: NewCartCommandCache(cacheStore),
	}
}
