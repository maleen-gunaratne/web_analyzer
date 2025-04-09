package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"

	"web-analyzer/internal/analysis"
	"web-analyzer/pkg/metrics"
)

func SetupRouter() *gin.Engine {

	setupLogger()

	router := gin.New()
	router.Use(
		gin.Recovery(),
		requestIDMiddleware(),
		loggerMiddleware(),
		configureCORS(),
	)
	metrics.InitMetrics()
	registerRoutes(router)

	return router
}

func RunServer(router *gin.Engine, addr string) {

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		slog.Info("starting server", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exited")
}

func registerRoutes(r *gin.Engine) {
	// Health and metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/health", healthCheckHandler)

	// Backward compatibility
	r.GET("/url_analyze", analysis.HandleAnalyze)

	// API v1
	api := r.Group("/api/v1")
	{
		api.GET("/analyze", analysis.HandleAnalyze)
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		slog.Info("request processed", "method", c.Request.Method, "path", c.Request.URL.Path,
			"status", c.Writer.Status(), "latency", latency, "requestID", c.GetString("requestID"), "clientIP", c.ClientIP(),
		)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func setupLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func configureCORS() gin.HandlerFunc {

	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
