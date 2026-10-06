// Package graphtest exposes a test-friendly GraphQL HTTP handler for the
// apigateway module.  External test packages import this instead of internal/.
package graphtest

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	graph "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/handler"
	graphqlmapper "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/mapper"
	"github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/middlewares"
	authpb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	pbbanner "github.com/MamangRust/monolith-graphql-ecommerce-pb/banner"
	pbcart "github.com/MamangRust/monolith-graphql-ecommerce-pb/cart"
	pbcategory "github.com/MamangRust/monolith-graphql-ecommerce-pb/category"
	pbmerchant "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant"
	pbmerchantaward "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_award"
	pbmerchantbusiness "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_business"
	pbmerchantdetail "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_detail"
	pbmerchantpolicy "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_policy"
	pbmsl "github.com/MamangRust/monolith-graphql-ecommerce-pb/merchant_social_link"
	pborder "github.com/MamangRust/monolith-graphql-ecommerce-pb/order"
	pborderitem "github.com/MamangRust/monolith-graphql-ecommerce-pb/order_item"
	pbproduct "github.com/MamangRust/monolith-graphql-ecommerce-pb/product"
	pbreview "github.com/MamangRust/monolith-graphql-ecommerce-pb/review"
	pbreviewdetail "github.com/MamangRust/monolith-graphql-ecommerce-pb/review_detail"
	pbrole "github.com/MamangRust/monolith-graphql-ecommerce-pb/role"
	pbshipping "github.com/MamangRust/monolith-graphql-ecommerce-pb/shipping_address"
	pbslider "github.com/MamangRust/monolith-graphql-ecommerce-pb/slider"
	pbtransaction "github.com/MamangRust/monolith-graphql-ecommerce-pb/transaction"
	pbuser "github.com/MamangRust/monolith-graphql-ecommerce-pb/user"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	"github.com/vektah/gqlparser/v2/ast"
	"google.golang.org/grpc"
)

// NewTestHandler builds a full gqlgen HTTP handler wired to the given gRPC
// connections and returns it as a plain http.Handler suitable for httptest.
func NewTestHandler(
	conns map[string]*grpc.ClientConn,
	cacheStore *cache.CacheStore,
	log logger.LoggerInterface,
) http.Handler {
	return newTestHandler(conns, cacheStore, log)
}

// NewAuthenticatedTestHandler behaves like NewTestHandler but wraps the
// GraphQL handler with the apigateway's production auth middleware, so
// protected operations (e.g. getMe) resolve the caller from the
// Authorization: Bearer header. Public operations (login/register/refresh)
// still pass through without a token.
func NewAuthenticatedTestHandler(
	conns map[string]*grpc.ClientConn,
	cacheStore *cache.CacheStore,
	log logger.LoggerInterface,
	tm auth.TokenManager,
) http.Handler {
	return middlewares.AuthMiddleware(tm, log)(newTestHandler(conns, cacheStore, log))
}

func newTestHandler(
	conns map[string]*grpc.ClientConn,
	cacheStore *cache.CacheStore,
	log logger.LoggerInterface,
) http.Handler {
	grpcClients := buildGRPCClients(conns)
	mapper := graphqlmapper.NewGraphqlMapper()
	imageUpload := upload_image.NewImageUpload(log)

	resolver := graph.NewResolver(&graph.Deps{
		Clients:     grpcClients,
		Logger:      log,
		Mapping:     mapper,
		Cache:       cacheStore,
		ImageUpload: imageUpload,
	})

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func buildGRPCClients(conns map[string]*grpc.ClientConn) *graph.GRPCClients {
	get := func(key string) *grpc.ClientConn {
		if c, ok := conns[key]; ok && c != nil {
			return c
		}
		return nil
	}

	return &graph.GRPCClients{
		AuthClient:                       authpb.NewAuthServiceClient(get("auth")),
		RoleCommandClient:                pbrole.NewRoleCommandServiceClient(get("role")),
		RoleQueryClient:                  pbrole.NewRoleQueryServiceClient(get("role")),
		UserCommandClient:                pbuser.NewUserCommandServiceClient(get("user")),
		UserQueryClient:                  pbuser.NewUserQueryServiceClient(get("user")),
		BannerCommandClient:              pbbanner.NewBannerCommandServiceClient(get("banner")),
		BannerQueryClient:                pbbanner.NewBannerQueryServiceClient(get("banner")),
		CartCommandClient:                pbcart.NewCartCommandServiceClient(get("cart")),
		CartQueryClient:                  pbcart.NewCartQueryServiceClient(get("cart")),
		CategoryCommandClient:            pbcategory.NewCategoryCommandServiceClient(get("category")),
		CategoryQueryClient:              pbcategory.NewCategoryQueryServiceClient(get("category")),
		CategoryStatsClient:              pbcategory.NewCategoryStatsServiceClient(get("category")),
		CategoryStatsByMerchantClient:    pbcategory.NewCategoryStatsByMerchantServiceClient(get("category")),
		CategoryStatsByIdClient:          pbcategory.NewCategoryStatsByIdServiceClient(get("category")),
		MerchantCommandClient:            pbmerchant.NewMerchantCommandServiceClient(get("merchant")),
		MerchantQueryClient:              pbmerchant.NewMerchantQueryServiceClient(get("merchant")),
		MerchantAwardCommandClient:       pbmerchantaward.NewMerchantAwardCommandServiceClient(get("merchant_award")),
		MerchantAwardQueryClient:         pbmerchantaward.NewMerchantAwardQueryServiceClient(get("merchant_award")),
		MerchantBusinessCommandClient:    pbmerchantbusiness.NewMerchantBusinessCommandServiceClient(get("merchant_business")),
		MerchantBusinessQueryClient:      pbmerchantbusiness.NewMerchantBusinessQueryServiceClient(get("merchant_business")),
		MerchantDetailCommandClient:      pbmerchantdetail.NewMerchantDetailCommandServiceClient(get("merchant_detail")),
		MerchantDetailQueryClient:        pbmerchantdetail.NewMerchantDetailQueryServiceClient(get("merchant_detail")),
		MerchantPolicyCommandClient:      pbmerchantpolicy.NewMerchantPolicyCommandServiceClient(get("merchant_policy")),
		MerchantPolicyQueryClient:        pbmerchantpolicy.NewMerchantPolicyQueryServiceClient(get("merchant_policy")),
		MerchantSocialLinkClient:         pbmsl.NewMerchantSocialCommandServiceClient(get("merchant")),
		OrderCommandClient:               pborder.NewOrderCommandServiceClient(get("order")),
		OrderQueryClient:                 pborder.NewOrderQueryServiceClient(get("order")),
		OrderStatsClient:                 pborder.NewOrderStatsServiceClient(get("order")),
		OrderStatsByMerchantClient:       pborder.NewOrderStatsByMerchantServiceClient(get("order")),
		OrderItemCommandClient:           pborderitem.NewOrderItemCommandServiceClient(get("order-item")),
		OrderItemQueryClient:             pborderitem.NewOrderItemQueryServiceClient(get("order-item")),
		ProductCommandClient:             pbproduct.NewProductCommandServiceClient(get("product")),
		ProductQueryClient:               pbproduct.NewProductQueryServiceClient(get("product")),
		ReviewCommandClient:              pbreview.NewReviewCommandServiceClient(get("review")),
		ReviewQueryClient:                pbreview.NewReviewQueryServiceClient(get("review")),
		ReviewDetailCommandClient:        pbreviewdetail.NewReviewDetailCommandServiceClient(get("review-detail")),
		ReviewDetailQueryClient:          pbreviewdetail.NewReviewDetailQueryServiceClient(get("review-detail")),
		ShippingCommandClient:            pbshipping.NewShippingCommandServiceClient(get("shipping-address")),
		ShippingQueryClient:              pbshipping.NewShippingQueryServiceClient(get("shipping-address")),
		SliderCommandClient:              pbslider.NewSliderCommandServiceClient(get("slider")),
		SliderQueryClient:                pbslider.NewSliderQueryServiceClient(get("slider")),
		TransactionCommandClient:         pbtransaction.NewTransactionCommandServiceClient(get("transaction")),
		TransactionQueryClient:           pbtransaction.NewTransactionQueryServiceClient(get("transaction")),
		TransactionStatsClient:           pbtransaction.NewTransactionStatsServiceClient(get("transaction")),
		TransactionStatsByMerchantClient: pbtransaction.NewTransactionStatsByMerchantServiceClient(get("transaction")),
	}
}
