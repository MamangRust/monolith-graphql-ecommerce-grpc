package apps

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/kafka"
	mencache "github.com/MamangRust/monolith-graphql-ecommerce-apigateway/internal/redis"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
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
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/dotenv"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-graphql-ecommerce-pkg/otel"
	redisclient "github.com/MamangRust/monolith-graphql-ecommerce-pkg/redis"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-graphql-ecommerce-pkg/upload_image"
	sharedcache "github.com/MamangRust/monolith-graphql-ecommerce-shared/cache"
	sharedobservability "github.com/MamangRust/monolith-graphql-ecommerce-shared/observability"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceAddresses struct {
	Auth             string
	Role             string
	User             string
	Category         string
	Merchant         string
	OrderItem        string
	Order            string
	Product          string
	Transaction      string
	Cart             string
	Review           string
	Slider           string
	Shipping         string
	Banner           string
	MerchantAward    string
	MerchantBusiness string
	MerchantDetail   string
	MerchantPolicy   string
	ReviewDetail     string
}

func loadServiceAddresses() *ServiceAddresses {
	return &ServiceAddresses{
		Auth:             getEnvOrDefault("GRPC_AUTH_ADDR", "localhost:50051"),
		Role:             getEnvOrDefault("GRPC_ROLE_ADDR", "localhost:50052"),
		User:             getEnvOrDefault("GRPC_USER_ADDR", "localhost:50053"),
		Category:         getEnvOrDefault("GRPC_CATEGORY_ADDR", "localhost:50054"),
		Merchant:         getEnvOrDefault("GRPC_MERCHANT_ADDR", "localhost:50055"),
		OrderItem:        getEnvOrDefault("GRPC_ORDER_ITEM_ADDR", "localhost:50056"),
		Order:            getEnvOrDefault("GRPC_ORDER_ADDR", "localhost:50057"),
		Product:          getEnvOrDefault("GRPC_PRODUCT_ADDR", "localhost:50058"),
		Transaction:      getEnvOrDefault("GRPC_TRANSACTION_ADDR", "localhost:50059"),
		Cart:             getEnvOrDefault("GRPC_CART_ADDR", "localhost:50060"),
		Review:           getEnvOrDefault("GRPC_REVIEW_ADDR", "localhost:50061"),
		Slider:           getEnvOrDefault("GRPC_SLIDER_ADDR", "localhost:50062"),
		Shipping:         getEnvOrDefault("GRPC_SHIPPING_ADDRESS_ADDR", "localhost:50063"),
		Banner:           getEnvOrDefault("GRPC_BANNER_ADDR", "localhost:50064"),
		MerchantAward:    getEnvOrDefault("GRPC_MERCHANT_AWARD_ADDR", "localhost:50065"),
		MerchantBusiness: getEnvOrDefault("GRPC_MERCHANT_BUSINESS_ADDR", "localhost:50066"),
		MerchantDetail:   getEnvOrDefault("GRPC_MERCHANT_DETAIL_ADDR", "localhost:50067"),
		MerchantPolicy:   getEnvOrDefault("GRPC_MERCHANT_POLICY_ADDR", "localhost:50068"),
		ReviewDetail:     getEnvOrDefault("GRPC_REVIEW_DETAIL_ADDR", "localhost:50069"),
	}
}

func createServiceConnections(addresses *ServiceAddresses, logger logger.LoggerInterface) (*graph.ServiceConnections, error) {
	var connections graph.ServiceConnections

	conns := map[string]*string{
		"Auth":             &addresses.Auth,
		"Role":             &addresses.Role,
		"User":             &addresses.User,
		"Category":         &addresses.Category,
		"Merchant":         &addresses.Merchant,
		"OrderItem":        &addresses.OrderItem,
		"Order":            &addresses.Order,
		"Product":          &addresses.Product,
		"Transaction":      &addresses.Transaction,
		"Cart":             &addresses.Cart,
		"Review":           &addresses.Review,
		"Slider":           &addresses.Slider,
		"Shipping":         &addresses.Shipping,
		"Banner":           &addresses.Banner,
		"MerchantAward":    &addresses.MerchantAward,
		"MerchantBusiness": &addresses.MerchantBusiness,
		"MerchantDetail":   &addresses.MerchantDetail,
		"MerchantPolicy":   &addresses.MerchantPolicy,
		"ReviewDetail":     &addresses.ReviewDetail,
	}

	for name, addr := range conns {
		conn, err := createConnection(*addr, name, logger)
		if err != nil {
			return nil, err
		}

		switch name {
		case "Auth":
			connections.AuthClient = conn
		case "Role":
			connections.RoleClient = conn
		case "User":
			connections.UserClient = conn
		case "Category":
			connections.CategoryClient = conn
		case "Merchant":
			connections.MerchantClient = conn
			connections.MerchantSocialLinkClient = conn
		case "OrderItem":
			connections.OrderItemClient = conn
		case "Order":
			connections.OrderClient = conn
		case "Product":
			connections.ProductClient = conn
		case "Transaction":
			connections.TransactionClient = conn
		case "Cart":
			connections.CartClient = conn
		case "Review":
			connections.ReviewClient = conn
		case "Slider":
			connections.SliderClient = conn
		case "Shipping":
			connections.ShippingClient = conn
		case "Banner":
			connections.BannerClient = conn
		case "MerchantAward":
			connections.MerchantAwardClient = conn
		case "MerchantBusiness":
			connections.MerchantBusinessClient = conn
		case "MerchantDetail":
			connections.MerchantDetailClient = conn
		case "MerchantPolicy":
			connections.MerchantPolicyClient = conn
		case "ReviewDetail":
			connections.ReviewDetailClient = conn
		}
	}

	return &connections, nil
}

func createConnection(address, serviceName string, logger logger.LoggerInterface) (*grpc.ClientConn, error) {
	logger.Info(fmt.Sprintf("Connecting to %s service at %s", serviceName, address))

	guard := resilience.NewDependencyGuardInterceptor(logger)

	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(guard.UnaryInterceptor()),
	)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to %s service", serviceName), zap.Error(err))
		return nil, err
	}
	return conn, nil
}

func closeConnections(conns *graph.ServiceConnections, log logger.LoggerInterface) {
	for name, conn := range map[string]*grpc.ClientConn{
		"Auth":               conns.AuthClient,
		"Role":               conns.RoleClient,
		"User":               conns.UserClient,
		"Category":           conns.CategoryClient,
		"Merchant":           conns.MerchantClient,
		"OrderItem":          conns.OrderItemClient,
		"Order":              conns.OrderClient,
		"Product":            conns.ProductClient,
		"Transaction":        conns.TransactionClient,
		"Cart":               conns.CartClient,
		"Review":             conns.ReviewClient,
		"Slider":             conns.SliderClient,
		"Shipping":           conns.ShippingClient,
		"Banner":             conns.BannerClient,
		"MerchantAward":      conns.MerchantAwardClient,
		"MerchantBusiness":   conns.MerchantBusinessClient,
		"MerchantDetail":     conns.MerchantDetailClient,
		"MerchantPolicy":     conns.MerchantPolicyClient,
		"ReviewDetail":       conns.ReviewDetailClient,
		"MerchantSocialLink": conns.MerchantSocialLinkClient,
	} {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Error(fmt.Sprintf("Failed to close %s connection", name), zap.Error(err))
			}
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

type Client struct {
	Logger logger.LoggerInterface
}

func RunClient() (*Client, func(), error) {
	flag.Parse()

	addresses := loadServiceAddresses()

	if err := dotenv.Viper(); err != nil {
		fmt.Printf("Warning: Failed to load .env file: %v\n", err)
	}

	ctx := context.Background()

	telemetry := otel_pkg.NewTelemetry(otel_pkg.Config{
		ServiceName:          "apigateway",
		ServiceVersion:       "1.0.0",
		Environment:          "development",
		Endpoint:             getEnvOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
		Insecure:             true,
		EnableRuntimeMetrics: true,
	})
	if err := telemetry.Init(ctx); err != nil {
		fmt.Printf("Warning: Failed to initialize telemetry: %v\n", err)
	}

	log, err := logger.NewLogger("apigateway", telemetry.GetLogger())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create logger: %w", err)
	}

	log.Debug("Creating gRPC connections...")
	conns, err := createServiceConnections(addresses, log)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect services: %w", err)
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		log.Fatal("Failed to create token manager", zap.Error(err))
	}

	myredis := redisclient.NewRedisClient(&redisclient.Config{
		Host:         viper.GetString("REDIS_HOST"),
		Port:         viper.GetString("REDIS_PORT"),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB_APIGATEWAY"),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	if err := myredis.Client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to ping redis", zap.Error(err))
	}

	cacheMetrics, err := sharedobservability.NewCacheMetrics("apigateway")
	if err != nil {
		log.Error("Failed to initialize cache metrics for apigateway cache store", zap.Error(err))
	}
	store := sharedcache.NewCacheStore(myredis.Client, log, cacheMetrics)

	imageUpload := upload_image.NewImageUpload(log)

	graphqlMapper := graphqlmapper.NewGraphqlMapper()

	grpcClients := &graph.GRPCClients{
		AuthClient:                       authpb.NewAuthServiceClient(conns.AuthClient),
		RoleCommandClient:                pbrole.NewRoleCommandServiceClient(conns.RoleClient),
		RoleQueryClient:                  pbrole.NewRoleQueryServiceClient(conns.RoleClient),
		UserCommandClient:                pbuser.NewUserCommandServiceClient(conns.UserClient),
		UserQueryClient:                  pbuser.NewUserQueryServiceClient(conns.UserClient),
		BannerCommandClient:              pbbanner.NewBannerCommandServiceClient(conns.BannerClient),
		BannerQueryClient:                pbbanner.NewBannerQueryServiceClient(conns.BannerClient),
		CartCommandClient:                pbcart.NewCartCommandServiceClient(conns.CartClient),
		CartQueryClient:                  pbcart.NewCartQueryServiceClient(conns.CartClient),
		CategoryCommandClient:            pbcategory.NewCategoryCommandServiceClient(conns.CategoryClient),
		CategoryQueryClient:              pbcategory.NewCategoryQueryServiceClient(conns.CategoryClient),
		CategoryStatsClient:              pbcategory.NewCategoryStatsServiceClient(conns.CategoryClient),
		CategoryStatsByMerchantClient:    pbcategory.NewCategoryStatsByMerchantServiceClient(conns.CategoryClient),
		CategoryStatsByIdClient:          pbcategory.NewCategoryStatsByIdServiceClient(conns.CategoryClient),
		MerchantCommandClient:            pbmerchant.NewMerchantCommandServiceClient(conns.MerchantClient),
		MerchantQueryClient:              pbmerchant.NewMerchantQueryServiceClient(conns.MerchantClient),
		MerchantAwardCommandClient:       pbmerchantaward.NewMerchantAwardCommandServiceClient(conns.MerchantAwardClient),
		MerchantAwardQueryClient:         pbmerchantaward.NewMerchantAwardQueryServiceClient(conns.MerchantAwardClient),
		MerchantBusinessCommandClient:    pbmerchantbusiness.NewMerchantBusinessCommandServiceClient(conns.MerchantBusinessClient),
		MerchantBusinessQueryClient:      pbmerchantbusiness.NewMerchantBusinessQueryServiceClient(conns.MerchantBusinessClient),
		MerchantDetailCommandClient:      pbmerchantdetail.NewMerchantDetailCommandServiceClient(conns.MerchantDetailClient),
		MerchantDetailQueryClient:        pbmerchantdetail.NewMerchantDetailQueryServiceClient(conns.MerchantDetailClient),
		MerchantPolicyCommandClient:      pbmerchantpolicy.NewMerchantPolicyCommandServiceClient(conns.MerchantPolicyClient),
		MerchantPolicyQueryClient:        pbmerchantpolicy.NewMerchantPolicyQueryServiceClient(conns.MerchantPolicyClient),
		MerchantSocialLinkClient:         pbmsl.NewMerchantSocialCommandServiceClient(conns.MerchantSocialLinkClient),
		OrderCommandClient:               pborder.NewOrderCommandServiceClient(conns.OrderClient),
		OrderQueryClient:                 pborder.NewOrderQueryServiceClient(conns.OrderClient),
		OrderStatsClient:                 pborder.NewOrderStatsServiceClient(conns.OrderClient),
		OrderStatsByMerchantClient:       pborder.NewOrderStatsByMerchantServiceClient(conns.OrderClient),
		OrderItemCommandClient:           pborderitem.NewOrderItemCommandServiceClient(conns.OrderItemClient),
		OrderItemQueryClient:             pborderitem.NewOrderItemQueryServiceClient(conns.OrderItemClient),
		ProductCommandClient:             pbproduct.NewProductCommandServiceClient(conns.ProductClient),
		ProductQueryClient:               pbproduct.NewProductQueryServiceClient(conns.ProductClient),
		ReviewCommandClient:              pbreview.NewReviewCommandServiceClient(conns.ReviewClient),
		ReviewQueryClient:                pbreview.NewReviewQueryServiceClient(conns.ReviewClient),
		ReviewDetailCommandClient:        pbreviewdetail.NewReviewDetailCommandServiceClient(conns.ReviewDetailClient),
		ReviewDetailQueryClient:          pbreviewdetail.NewReviewDetailQueryServiceClient(conns.ReviewDetailClient),
		ShippingCommandClient:            pbshipping.NewShippingCommandServiceClient(conns.ShippingClient),
		ShippingQueryClient:              pbshipping.NewShippingQueryServiceClient(conns.ShippingClient),
		SliderCommandClient:              pbslider.NewSliderCommandServiceClient(conns.SliderClient),
		SliderQueryClient:                pbslider.NewSliderQueryServiceClient(conns.SliderClient),
		TransactionCommandClient:         pbtransaction.NewTransactionCommandServiceClient(conns.TransactionClient),
		TransactionQueryClient:           pbtransaction.NewTransactionQueryServiceClient(conns.TransactionClient),
		TransactionStatsClient:           pbtransaction.NewTransactionStatsServiceClient(conns.TransactionClient),
		TransactionStatsByMerchantClient: pbtransaction.NewTransactionStatsByMerchantServiceClient(conns.TransactionClient),
	}

	myKafka := kafka.NewKafka(log, []string{os.Getenv("KAFKA_BROKERS")})

	mencache := mencache.NewCacheApiGateway(&mencache.Deps{
		Redis:  myredis.Client,
		Logger: log,
	})

	resolver := graph.NewResolver(&graph.Deps{
		Clients:     grpcClients,
		Logger:      log,
		Mapping:     graphqlMapper,
		Cache:       store,
		ImageUpload: imageUpload,
		Kafka:    myKafka,
		Mencache: mencache,
	})

	graphqlServer := setupGraphql(tokenManager, resolver, log, myredis.Client)

	go func() {
		log.Info(fmt.Sprintf("🚀 Starting GraphQL server on %s", graphqlServer.Addr))
		if err := graphqlServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("GraphQL server error", zap.Error(err))
		}
	}()

	metricsServer := &http.Server{
		Addr:              ":8091",
		Handler:           promhttp.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("Starting Prometheus metrics server on :8091")
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Metrics server error", zap.Error(err))
		}
	}()

	shutdown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Info("Shutting down GraphQL API Gateway...")

		if err := graphqlServer.Shutdown(ctx); err != nil {
			log.Error("GraphQL server shutdown failed", zap.Error(err))
		}

		if err := metricsServer.Shutdown(ctx); err != nil {
			log.Error("Metrics server shutdown failed", zap.Error(err))
		}

		closeConnections(conns, log)

		if err := telemetry.Shutdown(ctx); err != nil {
			log.Error("Telemetry shutdown failed", zap.Error(err))
		}

		log.Info("Shutdown complete ✅")
	}

	return &Client{
		Logger: log,
	}, shutdown, nil
}

func setupGraphql(token auth.TokenManager, resolver *graph.Resolver, logger logger.LoggerInterface, rdb *redis.Client) *http.Server {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
		Directives: graph.DirectiveRoot{
			HasRole: middlewares.HasRole(resolver.RoleGraphql.Permission),
		},
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

	port := getEnvOrDefault("CLIENT_PORT", "5000")

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)

	// k8s liveness/readiness probes target these paths; they must stay
	// registered on the GraphQL server port.
	r.Get("/health/live", createLivenessHandler())
	r.Get("/health/ready", createReadinessHandler(logger, rdb))

	r.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(token, logger))
		r.Handle("/query", srv)
	})

	logger.Info("GraphQL Playground running",
		zap.String("url", "http://localhost:"+port),
		zap.String("endpoint", "/query"),
	)

	return &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
}

func createLivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	}
}

func createReadinessHandler(logger logger.LoggerInterface, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err != nil {
			logger.Error("Readiness check failed: redis unreachable", zap.Error(err))
			respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"status": "not_ready",
				"deps":   map[string]string{"redis": "down"},
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "ready",
			"deps":   map[string]string{"redis": "up"},
			"time":   time.Now().UTC(),
		})
	}
}

func respondJSON(w http.ResponseWriter, code int, payload map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
