package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	"avito_service/internal/config"
	"avito_service/internal/storage/postgres"
	"avito_service/utils"
)

var (
	cfgFlag string
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	flag.StringVar(&cfgFlag, "path", "", "config-path")
	flag.Parse()
	cfg := config.LoadConfig(cfgFlag)
	logger := utils.ConfigureLogger(cfg.Env)
	logger.Info("config and logger successfully configured!")
	storage, err := postgres.NewStorage()
	if err != nil {
		logger.Error("failed to initialize db", err)
		os.Exit(1)
	}
	_ = storage
	logger.Info("database is up!")

	router := chi.NewRouter()
	utils.ConfigureRouter(router, logger, storage)

	server := &http.Server{
		Addr:        cfg.Address,
		Handler:     router,
		ReadTimeout: cfg.TimeOut,
		IdleTimeout: cfg.IdleTimeOut,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start server", err)
	}

}
