package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"web-analyzer/internal/server"
)

func main() {
	var (
		port        = flag.Int("port", 8080, "Port for the HTTP server")
		debugPort   = flag.Int("debug-port", 6060, "Debug server port for pprof")
		logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
		concurrency = flag.Int("concurrency", runtime.NumCPU(), "Maximum concurrency level")
	)
	flag.Parse()

	logger := setupLogger(*logLevel)
	slog.SetDefault(logger)

	runtime.GOMAXPROCS(*concurrency)

	go startDebugServer(*debugPort)

	router := server.SetupRouter()
	addr := fmt.Sprintf(":%d", *port)

	slog.Info("starting web analyzer ", "port", *port, "debug_port", *debugPort,
		"log_level", *logLevel, "concurrency", *concurrency, "go_version", runtime.Version())

	server.RunServer(router, addr)
}

func setupLogger(levelStr string) *slog.Logger {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

func startDebugServer(port int) {
	addr := fmt.Sprintf("localhost:%d", port)
	slog.Info("Starting pprof debug server", "address", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("pprof debug server failed", "error", err)
	}
}
