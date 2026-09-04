package review_cache

import "github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"

type reviewMencache struct {
	ReviewQueryCache
	ReviewCommandCache
}

type ReviewMencache interface {
	ReviewQueryCache
	ReviewCommandCache
}

func NewReviewMencache(cacheStore *cache.CacheStore) ReviewMencache {
	return &reviewMencache{
		ReviewQueryCache:   NewReviewQueryCache(cacheStore),
		ReviewCommandCache: NewReviewCommandCache(cacheStore),
	}
}
