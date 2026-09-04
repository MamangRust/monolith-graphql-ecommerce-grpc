package cart_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
)

const (
	cartAllCacheKey = "cart:all:page:%d:pageSize:%d:search:%s"
	ttlDefault      = 5 * time.Minute
)

type cartQueryCache struct {
	store *cache.CacheStore
}

func NewCartQueryCache(store *cache.CacheStore) *cartQueryCache {
	return &cartQueryCache{store: store}
}

func (c *cartQueryCache) GetCachedCarts(
	ctx context.Context,
	request *model.FindAllCartInput,
) (*model.APIResponsePaginationCart, bool) {

	var search string
	if request.Search != nil {
		search = *request.Search
	}

	key := fmt.Sprintf(
		cartAllCacheKey,
		request.Page,
		request.PageSize,
		search,
	)

	result, found := cache.GetFromCache[model.APIResponsePaginationCart](
		ctx,
		c.store,
		key,
	)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (c *cartQueryCache) SetCachedCarts(
	ctx context.Context,
	request *model.FindAllCartInput,
	resp *model.APIResponsePaginationCart,
) {
	if resp == nil {
		return
	}

	var search string
	if request.Search != nil {
		search = *request.Search
	}

	key := fmt.Sprintf(
		cartAllCacheKey,
		request.Page,
		request.PageSize,
		search,
	)

	cache.SetToCache(ctx, c.store, key, resp, ttlDefault)
}
