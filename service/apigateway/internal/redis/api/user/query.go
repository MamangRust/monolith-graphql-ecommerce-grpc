package user_cache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/model"
)

type userQueryCache struct {
	store *cache.CacheStore
}

func NewUserQueryCache(store *cache.CacheStore) UserQueryCache {
	return &userQueryCache{store: store}
}

func userCachePagination(req *model.FindAllUserInput) (int32, int32, string) {
	var page, pageSize int32
	var search string

	if req != nil {
		if req.Page != nil {
			page = *req.Page
		}
		if req.PageSize != nil {
			pageSize = *req.PageSize
		}
		if req.Search != nil {
			search = *req.Search
		}
	}

	return page, pageSize, search
}

func (s *userQueryCache) GetCachedUsersCache(ctx context.Context, req *model.FindAllUserInput) (*model.APIResponsePaginationUser, bool) {
	page, pageSize, search := userCachePagination(req)
	key := fmt.Sprintf(userAllCacheKey, page, pageSize, search)

	result, found := cache.GetFromCache[model.APIResponsePaginationUser](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *userQueryCache) SetCachedUsersCache(ctx context.Context, req *model.FindAllUserInput, data *model.APIResponsePaginationUser) {
	if data == nil {
		return
	}

	page, pageSize, search := userCachePagination(req)
	key := fmt.Sprintf(userAllCacheKey, page, pageSize, search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *userQueryCache) GetCachedUserActiveCache(ctx context.Context, req *model.FindAllUserInput) (*model.APIResponsePaginationUserDeleteAt, bool) {
	page, pageSize, search := userCachePagination(req)
	key := fmt.Sprintf(userActiveCacheKey, page, pageSize, search)

	result, found := cache.GetFromCache[model.APIResponsePaginationUserDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *userQueryCache) SetCachedUserActiveCache(ctx context.Context, req *model.FindAllUserInput, data *model.APIResponsePaginationUserDeleteAt) {
	if data == nil {
		return
	}

	page, pageSize, search := userCachePagination(req)
	key := fmt.Sprintf(userActiveCacheKey, page, pageSize, search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *userQueryCache) GetCachedUserTrashedCache(ctx context.Context, req *model.FindAllUserInput) (*model.APIResponsePaginationUserDeleteAt, bool) {
	page, pageSize, search := userCachePagination(req)
	key := fmt.Sprintf(userTrashedCacheKey, page, pageSize, search)

	result, found := cache.GetFromCache[model.APIResponsePaginationUserDeleteAt](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *userQueryCache) SetCachedUserTrashedCache(ctx context.Context, req *model.FindAllUserInput, data *model.APIResponsePaginationUserDeleteAt) {
	if data == nil {
		return
	}

	page, pageSize, search := userCachePagination(req)
	key := fmt.Sprintf(userTrashedCacheKey, page, pageSize, search)

	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}

func (s *userQueryCache) GetCachedUserCache(ctx context.Context, id int) (*model.APIResponseUserResponse, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := cache.GetFromCache[model.APIResponseUserResponse](ctx, s.store, key)

	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (s *userQueryCache) SetCachedUserCache(ctx context.Context, data *model.APIResponseUserResponse) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userByIdCacheKey, data.Data.ID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
