package utils

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"avito_service/internal/http_server/handlers/segment"
	"avito_service/internal/http_server/handlers/user"
	"avito_service/internal/storage/postgres"
)

func ConfigureRouter(mux *chi.Mux, logger *slog.Logger, storage *postgres.Storage) {
	mux.Use(middleware.RequestID)
	mux.Post("/create-segment", segment.SaveSegment(logger, storage))
	mux.Delete("/delete-segment", segment.DeleteSegment(logger, storage))
	mux.Post("/create-user", user.CreateUser(logger, storage))
	mux.Post("/addUser", user.AddUserToSeg(logger, storage))
	mux.Delete("/deleteFromUser", user.DeleteSegmFromUser(logger, storage))
}
func ConfigureLogger(env string) *slog.Logger {
	logger := &slog.Logger{}

	switch env {
	case EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}
