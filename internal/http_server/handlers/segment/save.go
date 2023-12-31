package segment

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"

	"avito_service/internal/http_server/handlers/segment/response"
	"avito_service/internal/storage"
	"avito_service/pkg"
)

type SegmentsSaver interface {
	CreateSegment(slug string) error
	IfSlugExists(slug string) error
}

func SaveSegment(logger *slog.Logger, saver SegmentsSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			middleware.GetReqID(r.Context()),
		)

		request := pkg.Segment{}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.Error("failed to decode request", err)
			response.WriteToJson(w, http.StatusInternalServerError, err)
			return
		}

		logger.Info("request decoded", slog.Any("request", request))

		if err := saver.IfSlugExists(request.Slug); err != nil {
			var slugErr storage.SegmentAlreadyExistsError
			if errors.As(err, &slugErr) {
				logger.Error("segment already exists", slugErr.Slug)
				response.WriteToJson(w, http.StatusBadRequest, "already exists")
				return
			}
		}

		if err := saver.CreateSegment(request.Slug); err != nil {
			logger.Error("failed to create user to db", err)
			response.WriteToJson(w, http.StatusInternalServerError, request)
			return
		}

		logger.Info("request successfully handled")
		response.WriteToJson(w, http.StatusCreated, request)
	}
}
