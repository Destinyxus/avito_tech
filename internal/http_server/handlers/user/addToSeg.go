package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"

	"avito_service/internal/http_server/handlers/segment/response"
	"avito_service/pkg"
)

type Adder interface {
	AddUserToSeg(list []string, id int) error
}

func AddUserToSeg(logger *slog.Logger, adder Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			middleware.GetReqID(r.Context()),
		)

		var listOfSegments pkg.SegmentToAdd

		if err := json.NewDecoder(r.Body).Decode(&listOfSegments); err != nil {
			logger.Error("failed to decode request", err)
			response.WriteToJson(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("request decoded", slog.Any("request", listOfSegments))

		if err := adder.AddUserToSeg(listOfSegments.Slug, listOfSegments.Id); err != nil {
			logger.Error("failed to save slug to db", err)
			response.WriteToJson(w, http.StatusInternalServerError, listOfSegments.Slug)
			return
		}

		logger.Info("request successfully handled")
		response.WriteToJson(w, http.StatusCreated, listOfSegments)
	}
}
