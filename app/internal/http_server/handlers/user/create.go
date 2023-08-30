package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"

	"avito_service/internal/http_server/handlers/segment/response"
	"avito_service/pkg"
)

type UserCreator interface {
	CreateUser(name string) error
}

func CreateUser(logger *slog.Logger, creator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			middleware.GetReqID(r.Context()),
		)

		user := pkg.User{}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Error("failed to decode request", err)
			response.WriteToJson(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("request decoded", slog.Any("request", user))

		if err := creator.CreateUser(user.Name); err != nil {
			logger.Error("failed to save slug to db", err)
			response.WriteToJson(w, http.StatusInternalServerError, user.Name)
			return
		}

		logger.Info("request successfully handled")
		response.WriteToJson(w, http.StatusCreated, user)
	}
}
