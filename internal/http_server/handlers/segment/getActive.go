package segment

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"

	"avito_service/internal/http_server/handlers/segment/response"
	"avito_service/internal/storage"
	"avito_service/pkg"
)

type GetterActive interface {
	GetActiveSegments(userid int) ([]string, error)
}

func GetActive(logger *slog.Logger, active GetterActive) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger = logger.With(
			middleware.GetReqID(request.Context()),
		)

		var activeRequest pkg.RequestActive

		if err := json.NewDecoder(request.Body).Decode(&activeRequest); err != nil {
			logger.Error("failed to decode request", err)
			response.WriteToJson(writer, http.StatusInternalServerError, err)
			return
		}

		logger.Info("request decoded", slog.Any("request", activeRequest))

		segments, err := active.GetActiveSegments(activeRequest.Id)
		if err != nil {
			var notFound storage.SegmentsNotFound
			if errors.As(err, &notFound) {
				logger.Error("active segments not found for this user", notFound)
				response.WriteToJson(writer, http.StatusNotFound, fmt.Sprintf("segments not found %v", segments))
				return
			} else {
				logger.Error("unexpected error", err)
				response.WriteToJson(writer, http.StatusInternalServerError, "")
				return
			}
		}
		logger.Info("active segments was successfully found!", segments)

		logger.Info("request successfully handled")
		response.WriteToJson(writer, http.StatusOK, segments)
	}
}
