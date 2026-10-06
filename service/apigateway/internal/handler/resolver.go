package graph

import (
	errorstd "errors"
	"fmt"
	"time"

	graphql "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/mapper"
	auth_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/auth"
	banner_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/banner"
	cart_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/cart"
	category_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/category"
	merchant_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/merchant"
	merchantawards_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/merchant_awards"
	merchantbusiness_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/merchant_business"
	merchantdetail_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/merchant_detail"
	merchantpolicies_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/merchant_policies"
	order_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/order"
	orderitem_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/order_item"
	product_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/product"
	review_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/review"
	reviewdetail_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/review_detail"
	role_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/role"
	shippingaddress_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/shipping_address"
	slider_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/slider"
	transaction_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/transaction"
	user_cache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis/api/user"
	authpb "github.com/MamangRust/monolith-graphql-ecommerce-pb"
	rolepermission "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/permission/role"
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
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/kafka"
	mencache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ServiceConnections struct {
	AuthClient               *grpc.ClientConn
	RoleClient               *grpc.ClientConn
	UserClient               *grpc.ClientConn
	CategoryClient           *grpc.ClientConn
	MerchantClient           *grpc.ClientConn
	OrderItemClient          *grpc.ClientConn
	OrderClient              *grpc.ClientConn
	ProductClient            *grpc.ClientConn
	TransactionClient        *grpc.ClientConn
	CartClient               *grpc.ClientConn
	ReviewClient             *grpc.ClientConn
	SliderClient             *grpc.ClientConn
	ShippingClient           *grpc.ClientConn
	BannerClient             *grpc.ClientConn
	MerchantAwardClient      *grpc.ClientConn
	MerchantBusinessClient   *grpc.ClientConn
	MerchantDetailClient     *grpc.ClientConn
	MerchantPolicyClient     *grpc.ClientConn
	ReviewDetailClient       *grpc.ClientConn
	MerchantSocialLinkClient *grpc.ClientConn
}

type Resolver struct {
	AuthGraphql               *AuthHandleGraphql
	RoleGraphql               *RoleHandleGraphql
	UserGraphql               *UserHandleGraphql
	CartGraphql               *CartHandleGraphql
	BannerGraphql             *BannerHandleGraphql
	CategoryGraphql           *CategoryHandleGraphql
	MerchantGraphql           *MerchantHandleGraphql
	MerchantAwardGraphql      *MerchantAwardHandleGraphql
	MerchantBusinessGraphql   *MerchantBusinessHandleGraphql
	MerchantDetailGraphql     *MerchantDetailHandleGraphql
	MerchantPolicyGraphql     *MerchantPolicyHandleGraphql
	MerchantSocialLinkGraphql *MerchantSocialLinkHandleGraphql
	OrderGraphql              *OrderHandleGraphql
	OrderItemGraphql          *OrderItemHandleGraphql
	ProductGraphql            *ProductHandleGraphql
	ReviewGraphql             *ReviewHandleGraphql
	ReviewDetailGraphql       *ReviewDetailHandleGraphql
	ShippingAddressGraphql    *ShippingAddressHandleGraphql
	SliderGraphql             *SliderHandleGraphql
	TransactionGraphql        *TransactionHandleGraphql
	ResolverHandle            *resolverHandler
}

type GRPCClients struct {
	AuthClient                       authpb.AuthServiceClient
	RoleCommandClient                pbrole.RoleCommandServiceClient
	RoleQueryClient                  pbrole.RoleQueryServiceClient
	UserCommandClient                pbuser.UserCommandServiceClient
	UserQueryClient                  pbuser.UserQueryServiceClient
	BannerCommandClient              pbbanner.BannerCommandServiceClient
	BannerQueryClient                pbbanner.BannerQueryServiceClient
	CartCommandClient                pbcart.CartCommandServiceClient
	CartQueryClient                  pbcart.CartQueryServiceClient
	CategoryCommandClient            pbcategory.CategoryCommandServiceClient
	CategoryQueryClient              pbcategory.CategoryQueryServiceClient
	CategoryStatsClient              pbcategory.CategoryStatsServiceClient
	CategoryStatsByMerchantClient    pbcategory.CategoryStatsByMerchantServiceClient
	CategoryStatsByIdClient          pbcategory.CategoryStatsByIdServiceClient
	MerchantCommandClient            pbmerchant.MerchantCommandServiceClient
	MerchantQueryClient              pbmerchant.MerchantQueryServiceClient
	MerchantAwardCommandClient       pbmerchantaward.MerchantAwardCommandServiceClient
	MerchantAwardQueryClient         pbmerchantaward.MerchantAwardQueryServiceClient
	MerchantBusinessCommandClient    pbmerchantbusiness.MerchantBusinessCommandServiceClient
	MerchantBusinessQueryClient      pbmerchantbusiness.MerchantBusinessQueryServiceClient
	MerchantDetailCommandClient      pbmerchantdetail.MerchantDetailCommandServiceClient
	MerchantDetailQueryClient        pbmerchantdetail.MerchantDetailQueryServiceClient
	MerchantPolicyCommandClient      pbmerchantpolicy.MerchantPolicyCommandServiceClient
	MerchantPolicyQueryClient        pbmerchantpolicy.MerchantPolicyQueryServiceClient
	MerchantSocialLinkClient         pbmsl.MerchantSocialCommandServiceClient
	OrderCommandClient               pborder.OrderCommandServiceClient
	OrderQueryClient                 pborder.OrderQueryServiceClient
	OrderStatsClient                 pborder.OrderStatsServiceClient
	OrderStatsByMerchantClient       pborder.OrderStatsByMerchantServiceClient
	OrderItemCommandClient           pborderitem.OrderItemCommandServiceClient
	OrderItemQueryClient             pborderitem.OrderItemQueryServiceClient
	ProductCommandClient             pbproduct.ProductCommandServiceClient
	ProductQueryClient               pbproduct.ProductQueryServiceClient
	ReviewCommandClient              pbreview.ReviewCommandServiceClient
	ReviewQueryClient                pbreview.ReviewQueryServiceClient
	ReviewDetailCommandClient        pbreviewdetail.ReviewDetailCommandServiceClient
	ReviewDetailQueryClient          pbreviewdetail.ReviewDetailQueryServiceClient
	ShippingCommandClient            pbshipping.ShippingCommandServiceClient
	ShippingQueryClient              pbshipping.ShippingQueryServiceClient
	SliderCommandClient              pbslider.SliderCommandServiceClient
	SliderQueryClient                pbslider.SliderQueryServiceClient
	TransactionCommandClient         pbtransaction.TransactionCommandServiceClient
	TransactionQueryClient           pbtransaction.TransactionQueryServiceClient
	TransactionStatsClient           pbtransaction.TransactionStatsServiceClient
	TransactionStatsByMerchantClient pbtransaction.TransactionStatsByMerchantServiceClient
}

type Deps struct {
	Clients     *GRPCClients
	Logger      logger.LoggerInterface
	Mapping     *graphql.GraphqlMapper
	Cache       *cache.CacheStore
	ImageUpload upload_image.ImageUploads
	Kafka    *kafka.Kafka
	Mencache mencache.CacheApiGateway
}

func NewResolver(deps *Deps) *Resolver {
	obs, _ := observability.NewObservability(
		"graphql-client",
		deps.Logger,
	)

	resolver := NewResolverHandler(obs, deps.Logger)

	return &Resolver{
		AuthGraphql: &AuthHandleGraphql{
			AuthClient: deps.Clients.AuthClient,
			Mapping:    deps.Mapping.AuthGraphqlMapper,
			Logger:     deps.Logger,
			Cache:      auth_cache.NewMencache(deps.Cache),
		},
		RoleGraphql: &RoleHandleGraphql{
			RoleCommandClient: deps.Clients.RoleCommandClient,
			RoleQueryClient:   deps.Clients.RoleQueryClient,
			Mapping:           deps.Mapping.RoleGraphqlMapper,
			Logger:            deps.Logger,
			Cache:             role_cache.NewRoleMencache(deps.Cache),
			Permission:        rolepermission.NewRolePermission(deps.Kafka, "request-role", "response-role", 5*time.Second, deps.Logger, deps.Mencache),
		},
		UserGraphql: &UserHandleGraphql{
			UserCommandClient: deps.Clients.UserCommandClient,
			UserQueryClient:   deps.Clients.UserQueryClient,
			Mapping:           deps.Mapping.UserGraphqlMapper,
			Logger:            deps.Logger,
			Cache:             user_cache.NewUserMencache(deps.Cache),
		},
		BannerGraphql: &BannerHandleGraphql{
			BannerCommandClient: deps.Clients.BannerCommandClient,
			BannerQueryClient:   deps.Clients.BannerQueryClient,
			Mapping:             deps.Mapping.BannerGraphqlMapper,
			Logger:              deps.Logger,
			Cache:               banner_cache.NewBannerMencache(deps.Cache),
		},
		CartGraphql: &CartHandleGraphql{
			CartCommandClient: deps.Clients.CartCommandClient,
			CartQueryClient:   deps.Clients.CartQueryClient,
			Mapping:           deps.Mapping.CartGraphqlMapper,
			Logger:            deps.Logger,
			Cache:             cart_cache.NewCartMencache(deps.Cache),
		},
		CategoryGraphql: &CategoryHandleGraphql{
			CategoryCommandClient:         deps.Clients.CategoryCommandClient,
			CategoryQueryClient:           deps.Clients.CategoryQueryClient,
			CategoryStatsClient:           deps.Clients.CategoryStatsClient,
			CategoryStatsByMerchantClient: deps.Clients.CategoryStatsByMerchantClient,
			CategoryStatsByIdClient:       deps.Clients.CategoryStatsByIdClient,
			Mapping:                       deps.Mapping.CategoryGraphqlMapper,
			Logger:                        deps.Logger,
			Cache:                         category_cache.NewCategoryMencache(deps.Cache),
			UploadImage:                   deps.ImageUpload,
		},
		MerchantGraphql: &MerchantHandleGraphql{
			MerchantCommandClient: deps.Clients.MerchantCommandClient,
			MerchantQueryClient:   deps.Clients.MerchantQueryClient,
			Mapping:               deps.Mapping.MerchantGraphqlMapper,
			Logger:                deps.Logger,
			Cache:                 merchant_cache.NewMerchantMencache(deps.Cache),
		},
		MerchantAwardGraphql: &MerchantAwardHandleGraphql{
			MerchantAwardCommandClient: deps.Clients.MerchantAwardCommandClient,
			MerchantAwardQueryClient:   deps.Clients.MerchantAwardQueryClient,
			Mapping:                    deps.Mapping.MerchantAwardGraphqlMapper,
			Logger:                     deps.Logger,
			Cache:                      merchantawards_cache.NewMerchantAward(deps.Cache),
		},
		MerchantBusinessGraphql: &MerchantBusinessHandleGraphql{
			MerchantBusinessCommandClient: deps.Clients.MerchantBusinessCommandClient,
			MerchantBusinessQueryClient:   deps.Clients.MerchantBusinessQueryClient,
			Mapping:                       deps.Mapping.MerchantBusinessGraphqlMapper,
			Logger:                        deps.Logger,
			Cache:                         merchantbusiness_cache.NewMerchantBusinessMencache(deps.Cache),
		},
		MerchantDetailGraphql: &MerchantDetailHandleGraphql{
			MerchantDetailCommandClient: deps.Clients.MerchantDetailCommandClient,
			MerchantDetailQueryClient:   deps.Clients.MerchantDetailQueryClient,
			Mapping:                     deps.Mapping.MerchantDetailGraphqlMapper,
			UploadImage:                 deps.ImageUpload,
			Logger:                      deps.Logger,
			Cache:                       merchantdetail_cache.NewMerchantDetailMencache(deps.Cache),
		},
		MerchantPolicyGraphql: &MerchantPolicyHandleGraphql{
			MerchantPolicyCommandClient: deps.Clients.MerchantPolicyCommandClient,
			MerchantPolicyQueryClient:   deps.Clients.MerchantPolicyQueryClient,
			Mapping:                     deps.Mapping.MerchantPolicyGraphqlMapper,
			Logger:                      deps.Logger,
			Cache:                       merchantpolicies_cache.NewMerchantPoliciesMencache(deps.Cache),
		},
		MerchantSocialLinkGraphql: &MerchantSocialLinkHandleGraphql{
			MerchantSocialLinkClient: deps.Clients.MerchantSocialLinkClient,
			Mapping:                  deps.Mapping.MerchantSocialLinkGraphqlMapper,
			Logger:                   deps.Logger,
		},
		OrderGraphql: &OrderHandleGraphql{
			OrderCommandClient:         deps.Clients.OrderCommandClient,
			OrderQueryClient:           deps.Clients.OrderQueryClient,
			OrderStatsClient:           deps.Clients.OrderStatsClient,
			OrderStatsByMerchantClient: deps.Clients.OrderStatsByMerchantClient,
			Mapping:                    deps.Mapping.OrderGraphqlMapper,
			Logger:                     deps.Logger,
			Cache:                      order_cache.OrderNewMencache(deps.Cache),
		},
		OrderItemGraphql: &OrderItemHandleGraphql{
			OrderItemCommandClient: deps.Clients.OrderItemCommandClient,
			OrderItemQueryClient:   deps.Clients.OrderItemQueryClient,
			Mapping:                deps.Mapping.OrderItemGraphqlMapper,
			Logger:                 deps.Logger,
			Cache:                  orderitem_cache.NewOrderItemMencache(deps.Cache),
		},
		ProductGraphql: &ProductHandleGraphql{
			ProductCommandClient: deps.Clients.ProductCommandClient,
			ProductQueryClient:   deps.Clients.ProductQueryClient,
			Mapping:              deps.Mapping.ProductGraphqlMapper,
			UploadImage:          deps.ImageUpload,
			Logger:               deps.Logger,
			Cache:                product_cache.NewProductMencache(deps.Cache),
		},
		ReviewGraphql: &ReviewHandleGraphql{
			ReviewCommandClient: deps.Clients.ReviewCommandClient,
			ReviewQueryClient:   deps.Clients.ReviewQueryClient,
			Mapping:             deps.Mapping.ReviewGraphqlMapper,
			Logger:              deps.Logger,
			Cache:               review_cache.NewReviewMencache(deps.Cache),
		},
		ReviewDetailGraphql: &ReviewDetailHandleGraphql{
			ReviewDetailCommandClient: deps.Clients.ReviewDetailCommandClient,
			ReviewDetailQueryClient:   deps.Clients.ReviewDetailQueryClient,
			Mapping:                   deps.Mapping.ReviewDetailGraphqlMapper,
			Logger:                    deps.Logger,
			Cache:                     reviewdetail_cache.NewReviewDetailMencache(deps.Cache),
		},
		ShippingAddressGraphql: &ShippingAddressHandleGraphql{
			ShippingCommandClient: deps.Clients.ShippingCommandClient,
			ShippingQueryClient:   deps.Clients.ShippingQueryClient,
			Mapping:               deps.Mapping.ShippingAddresGraphqlMapper,
			Logger:                deps.Logger,
			Cache:                 shippingaddress_cache.NewShippingAddressMencache(deps.Cache),
		},
		SliderGraphql: &SliderHandleGraphql{
			SliderCommandClient: deps.Clients.SliderCommandClient,
			SliderQueryClient:   deps.Clients.SliderQueryClient,
			Mapping:             deps.Mapping.SliderGraphqlMapper,
			UploadImage:         deps.ImageUpload,
			Logger:              deps.Logger,
			Cache:               slider_cache.NewSliderMencache(deps.Cache),
		},
		TransactionGraphql: &TransactionHandleGraphql{
			TransactionCommandClient:         deps.Clients.TransactionCommandClient,
			TransactionQueryClient:           deps.Clients.TransactionQueryClient,
			TransactionStatsClient:           deps.Clients.TransactionStatsClient,
			TransactionStatsByMerchantClient: deps.Clients.TransactionStatsByMerchantClient,
			Mapping:                          deps.Mapping.TransactionGraphqlMapper,
			Logger:                           deps.Logger,
			Cache:                            transaction_cache.NewTransactionMencache(deps.Cache),
		},
		ResolverHandle: resolver,
	}
}

type AuthHandleGraphql struct {
	AuthClient authpb.AuthServiceClient
	Mapping    graphql.AuthGraphqlMapper
	Logger     logger.LoggerInterface
	Cache      auth_cache.AuthMencache
}

type RoleHandleGraphql struct {
	RoleCommandClient pbrole.RoleCommandServiceClient
	RoleQueryClient   pbrole.RoleQueryServiceClient
	Mapping           graphql.RoleGraphqlMapper
	Logger            logger.LoggerInterface
	Cache             role_cache.RoleMencache
	Permission        rolepermission.RolePermission
}

type UserHandleGraphql struct {
	UserCommandClient pbuser.UserCommandServiceClient
	UserQueryClient   pbuser.UserQueryServiceClient
	Mapping           graphql.UserGraphqlMapper
	Logger            logger.LoggerInterface
	Cache             user_cache.UserMencache
}

type BannerHandleGraphql struct {
	BannerCommandClient pbbanner.BannerCommandServiceClient
	BannerQueryClient   pbbanner.BannerQueryServiceClient
	Mapping             graphql.BannerGraphqlMapper
	Logger              logger.LoggerInterface
	Cache               banner_cache.BannerMencache
}

type CartHandleGraphql struct {
	CartCommandClient pbcart.CartCommandServiceClient
	CartQueryClient   pbcart.CartQueryServiceClient
	Mapping           graphql.CartGraphqlMapper
	Logger            logger.LoggerInterface
	Cache             cart_cache.CartMencache
}

type CategoryHandleGraphql struct {
	CategoryCommandClient         pbcategory.CategoryCommandServiceClient
	CategoryQueryClient           pbcategory.CategoryQueryServiceClient
	CategoryStatsClient           pbcategory.CategoryStatsServiceClient
	CategoryStatsByMerchantClient pbcategory.CategoryStatsByMerchantServiceClient
	CategoryStatsByIdClient       pbcategory.CategoryStatsByIdServiceClient
	Mapping                       graphql.CategoryGraphqlMapper
	UploadImage                   upload_image.ImageUploads
	Logger                        logger.LoggerInterface
	Cache                         category_cache.CategoryMencache
}

type MerchantHandleGraphql struct {
	MerchantCommandClient pbmerchant.MerchantCommandServiceClient
	MerchantQueryClient   pbmerchant.MerchantQueryServiceClient
	Mapping               graphql.MerchantGraphqlMapper
	Logger                logger.LoggerInterface
	Cache                 merchant_cache.MerchantMencache
}

type MerchantAwardHandleGraphql struct {
	MerchantAwardCommandClient pbmerchantaward.MerchantAwardCommandServiceClient
	MerchantAwardQueryClient   pbmerchantaward.MerchantAwardQueryServiceClient
	Mapping                    graphql.MerchantAwardGraphqlMapper
	Logger                     logger.LoggerInterface
	Cache                      merchantawards_cache.MerchantAwardMencache
}

type MerchantBusinessHandleGraphql struct {
	MerchantBusinessCommandClient pbmerchantbusiness.MerchantBusinessCommandServiceClient
	MerchantBusinessQueryClient   pbmerchantbusiness.MerchantBusinessQueryServiceClient
	Mapping                       graphql.MerchantBusinessGraphqlMapper
	Logger                        logger.LoggerInterface
	Cache                         merchantbusiness_cache.MerchantBusinessMencache
}

type MerchantDetailHandleGraphql struct {
	MerchantDetailCommandClient pbmerchantdetail.MerchantDetailCommandServiceClient
	MerchantDetailQueryClient   pbmerchantdetail.MerchantDetailQueryServiceClient
	Mapping                     graphql.MerchantDetailGraphqlMapper
	UploadImage                 upload_image.ImageUploads
	Logger                      logger.LoggerInterface
	Cache                       merchantdetail_cache.MerchantDetailMencache
}

type MerchantPolicyHandleGraphql struct {
	MerchantPolicyCommandClient pbmerchantpolicy.MerchantPolicyCommandServiceClient
	MerchantPolicyQueryClient   pbmerchantpolicy.MerchantPolicyQueryServiceClient
	Mapping                     graphql.MerchantPolicyGraphqlMapper
	Logger                      logger.LoggerInterface
	Cache                       merchantpolicies_cache.MerchantPoliciesMencache
}

type MerchantSocialLinkHandleGraphql struct {
	MerchantSocialLinkClient pbmsl.MerchantSocialCommandServiceClient
	Mapping                  graphql.MerchantSocialLinkGraphqlMapper
	Logger                   logger.LoggerInterface
}

type OrderHandleGraphql struct {
	OrderCommandClient         pborder.OrderCommandServiceClient
	OrderQueryClient           pborder.OrderQueryServiceClient
	OrderStatsClient           pborder.OrderStatsServiceClient
	OrderStatsByMerchantClient pborder.OrderStatsByMerchantServiceClient
	Mapping                    graphql.OrderGraphqlMapper
	Logger                     logger.LoggerInterface
	Cache                      order_cache.OrderMencache
}

type OrderItemHandleGraphql struct {
	OrderItemCommandClient pborderitem.OrderItemCommandServiceClient
	OrderItemQueryClient   pborderitem.OrderItemQueryServiceClient
	Mapping                graphql.OrderItemGraphqlMapper
	Logger                 logger.LoggerInterface
	Cache                  orderitem_cache.OrderItemMencache
}

type ProductHandleGraphql struct {
	ProductCommandClient pbproduct.ProductCommandServiceClient
	ProductQueryClient   pbproduct.ProductQueryServiceClient
	Mapping              graphql.ProductGraphqlMapper
	UploadImage          upload_image.ImageUploads
	Logger               logger.LoggerInterface
	Cache                product_cache.ProductMencache
}

type ReviewHandleGraphql struct {
	ReviewCommandClient pbreview.ReviewCommandServiceClient
	ReviewQueryClient   pbreview.ReviewQueryServiceClient
	Mapping             graphql.ReviewGraphqlMapper
	Logger              logger.LoggerInterface
	Cache               review_cache.ReviewMencache
}

type ReviewDetailHandleGraphql struct {
	ReviewDetailCommandClient pbreviewdetail.ReviewDetailCommandServiceClient
	ReviewDetailQueryClient   pbreviewdetail.ReviewDetailQueryServiceClient
	Mapping                   graphql.ReviewDetailGraphqlMapper
	Logger                    logger.LoggerInterface
	Cache                     reviewdetail_cache.ReviewDetailMencache
}

type ShippingAddressHandleGraphql struct {
	ShippingCommandClient pbshipping.ShippingCommandServiceClient
	ShippingQueryClient   pbshipping.ShippingQueryServiceClient
	Mapping               graphql.ShippingAddresGraphqlMapper
	Logger                logger.LoggerInterface
	Cache                 shippingaddress_cache.ShippingAddressMencache
}

type SliderHandleGraphql struct {
	SliderCommandClient pbslider.SliderCommandServiceClient
	SliderQueryClient   pbslider.SliderQueryServiceClient
	Mapping             graphql.SliderGraphqlMapper
	UploadImage         upload_image.ImageUploads
	Logger              logger.LoggerInterface
	Cache               slider_cache.SliderMencache
}

type TransactionHandleGraphql struct {
	TransactionCommandClient         pbtransaction.TransactionCommandServiceClient
	TransactionQueryClient           pbtransaction.TransactionQueryServiceClient
	TransactionStatsClient           pbtransaction.TransactionStatsServiceClient
	TransactionStatsByMerchantClient pbtransaction.TransactionStatsByMerchantServiceClient
	Mapping                          graphql.TransactionGraphqlMapper
	Logger                           logger.LoggerInterface
	Cache                            transaction_cache.TransactionMencache
}

func (h *Resolver) handleGraphQLError(err error, operation string) *errors.AppError {
	if err == nil {
		return nil
	}

	var appErr *errors.AppError
	if errorstd.As(err, &appErr) {
		return appErr
	}

	return errors.NewInternalError(err).WithMessage("Failed to " + operation)
}

func (h *Resolver) parseValidationErrors(err error) []errors.ValidationError {
	var validationErrs []errors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, errors.ValidationError{
				Field:   fe.Field(),
				Message: h.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []errors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (h *Resolver) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
