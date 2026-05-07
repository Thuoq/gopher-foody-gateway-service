package http

import (
	"gopher-foody-gateway-service/internal/config"
	"gopher-foody-gateway-service/internal/presentation/http/handlers"
	"gopher-foody-gateway-service/internal/presentation/http/middleware"
	"gopher-foody-gateway-service/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(cfg *config.Config, logger *zap.Logger, jwtManager jwt.TokenManager) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Add basic middlewares
	r.Use(gin.Recovery())

	// Example health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Setup Proxy Handlers
	identityProxy, err := handlers.NewProxyHandler(cfg.Upstream.IdentityServiceURL)
	if err != nil {
		logger.Fatal("Failed to create identity proxy", zap.Error(err))
	}

	restaurantProxy, err := handlers.NewProxyHandler(cfg.Upstream.RestaurantServiceURL)
	if err != nil {
		logger.Fatal("Failed to create restaurant proxy", zap.Error(err))
	}

	orderProxy, err := handlers.NewProxyHandler(cfg.Upstream.OrderServiceURL)
	if err != nil {
		logger.Fatal("Failed to create order proxy", zap.Error(err))
	}

	api := r.Group("/api/v1")

	// Identity routes
	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", identityProxy.Proxy)
		auth.POST("/sign-in", identityProxy.Proxy)
		auth.POST("/refresh", identityProxy.Proxy)

		protectedAuth := auth.Group("")
		protectedAuth.Use(middleware.AuthMiddleware(jwtManager))
		{
			protectedAuth.POST("/logout", identityProxy.Proxy)
		}
	}

	// Restaurant routes
	restaurants := api.Group("/restaurants")
	{
		// Public
		restaurants.GET("", restaurantProxy.Proxy)
		restaurants.GET("/:id", restaurantProxy.Proxy)
		restaurants.GET("/:id/foods", restaurantProxy.Proxy)

		// Admin (Protected)
		admin := restaurants.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtManager))
		{
			admin.Any("/*path", restaurantProxy.Proxy)
		}
	}

	// Order routes
	orders := api.Group("/orders")
	orders.Use(middleware.AuthMiddleware(jwtManager))
	{
		orders.POST("", orderProxy.Proxy)
		orders.GET("", orderProxy.Proxy)
		orders.GET("/:id", orderProxy.Proxy)
	}

	return r
}
