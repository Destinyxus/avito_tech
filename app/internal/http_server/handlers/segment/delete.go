package segment

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"

	"avito_service/internal/http_server/handlers/segment/response"
	"avito_service/pkg"
)

type SegmentDeleter interface {
	DeleteSegment(slug string) error
}

func DeleteSegment(logger *slog.Logger, deleter SegmentDeleter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger = logger.With(
			middleware.GetReqID(request.Context()),
		)
		slug := pkg.Segment{}

		if err := json.NewDecoder(request.Body).Decode(&slug); err != nil {
			logger.Error("failed to decode request", err)
			response.WriteToJson(writer, http.StatusInternalServerError, err)
			return
		}

		logger.Info("request decoded", slog.Any("request", slug.Slug))

		if err := deleter.DeleteSegment(slug.Slug); err != nil {
			logger.Error(fmt.Sprintf("failed to delete slug %s", slug.Slug), err)
			response.WriteToJson(writer, http.StatusInternalServerError, request)
			return
		}

		logger.Info("request successfully handled")

		response.WriteToJson(writer, http.StatusOK, "segment was deleted")
	}

}
