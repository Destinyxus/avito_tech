package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"avito_service/internal/config"
	"avito_service/internal/service/configuration"
	"avito_service/internal/service/ttl"
	"avito_service/internal/storage/postgres"
)

var (
	cfgFlag string
)

func main() {
	flag.StringVar(&cfgFlag, "path", "", "config-path")
	flag.Parse()
	cfg := config.LoadConfig(cfgFlag)
	logger := configuration.ConfigureLogger(cfg.Env)
	logger.Info("config and logger successfully configured!")
	storage, err := postgres.NewStorage()
	if err != nil {
		logger.Error("failed to initialize db", err)
		os.Exit(1)
	}
	logger.Info("database is up!")

	go func() {
		if err := ttl.TTLChecker(logger, storage); err != nil {
			logger.Error("checker not launched..", err)
			os.Exit(1)
		}

	}()

	router := chi.NewRouter()
	configuration.ConfigureRouter(router, logger, storage)

	server := &http.Server{
		Addr:        cfg.Address,
		Handler:     router,
		ReadTimeout: cfg.TimeOut,
		IdleTimeout: cfg.IdleTimeOut,
	}

	logger.Info("Starting server...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start server", err)
	}

}
