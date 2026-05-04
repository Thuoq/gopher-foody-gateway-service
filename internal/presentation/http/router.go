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

	api := r.Group("/api/v1")

	// Public routes
	{
		api.POST("/auth/sign-up", identityProxy.Proxy)
		api.POST("/auth/sign-in", identityProxy.Proxy)
		api.POST("/auth/refresh", identityProxy.Proxy)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Identity protected routes
		protected.POST("/auth/logout", identityProxy.Proxy)

		// Other services
		protected.Any("/restaurants/*path", restaurantProxy.Proxy)
	}

	return r
}
