package main

import (
	"context"
	"fmt"
	"gopher-foody-gateway-service/internal/config"
	"gopher-foody-gateway-service/pkg/jwt"
	"gopher-foody-gateway-service/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Core dependencies
	container.Provide(config.LoadConfig)
	container.Provide(logger.NewLogger)
	container.Provide(jwt.NewManager)

	// Presentation
	container.Provide(httpRouter.NewRouter)

	return container
}

func main() {
	container := BuildContainer()

	err := container.Invoke(func(cfg *config.Config, log *zap.Logger, router *gin.Engine) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Start HTTP Server
		httpServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.App.HTTPPort),
			Handler: router,
		}

		go func() {
			log.Info("Starting Gateway HTTP Server", zap.Int("port", cfg.App.HTTPPort))
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal("Gateway Server failed to start", zap.Error(err))
			}
		}()

		// Wait for interrupt signal to gracefully shutdown
		<-ctx.Done()
		log.Info("Shutting down Gateway gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error("Gateway Server forced to shutdown", zap.Error(err))
		}

		log.Info("Gateway exited gracefully")
	})

	if err != nil {
		fmt.Printf("Error starting application: %v\n", err)
		os.Exit(1)
	}
}
