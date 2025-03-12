package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vit6556/ozon-internship-assignment/internal/app"
	"github.com/vit6556/ozon-internship-assignment/internal/config"
)

func main() {
	var serverConfig config.HTTPServerConfig
	config.LoadConfig(&serverConfig)

	echo, dbPool := app.InitServer(&serverConfig)

	go func() {
		if err := echo.Start(fmt.Sprintf(":%d", serverConfig.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := echo.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %s", err)
	}

	if dbPool != nil {
		log.Println("closing database connection...")
		dbPool.Close()
	}

	log.Println("server stopped gracefully")
}
