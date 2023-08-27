package user

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

type UserSegDeletor interface {
	DeleteSegmentsOfUser(list []string, userid int) error
	IfExists(userId int, slugList []string) error
}

func DeleteSegmFromUser(logger *slog.Logger, deletor UserSegDeletor) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger = logger.With(
			middleware.GetReqID(request.Context()),
		)

		var deleteRequest pkg.SegmentToAdd

		if err := json.NewDecoder(request.Body).Decode(&deleteRequest); err != nil {
			logger.Error("failed to decode request", err)
			response.WriteToJson(writer, http.StatusInternalServerError, err)
			return
		}

		logger.Info("request decoded", slog.Any("request", deleteRequest))

		if err := deletor.IfExists(deleteRequest.Id, deleteRequest.Slug); err != nil {
			var userError storage.UserNotExists
			var segmentError storage.SegmentNotExists
			if errors.As(err, &userError) {
				logger.Error("provided user not exists", userError)
				response.WriteToJson(writer, http.StatusConflict, fmt.Sprintf("user %d not exists", userError.Id))
				return
			} else if errors.As(err, &segmentError) {
				logger.Error("provided segment is not active", segmentError)
				response.WriteToJson(writer, http.StatusNotFound, fmt.Sprintf("slug %s not exists", segmentError.Slug))
				return
			} else {
				logger.Error("unexpected error from db")
				response.WriteToJson(writer, http.StatusInternalServerError, "")
				return
			}
		}

		if err := deletor.DeleteSegmentsOfUser(deleteRequest.Slug, deleteRequest.Id); err != nil {
			logger.Error("error", err)
			response.WriteToJson(writer, http.StatusInternalServerError, "")
			return
		}

		logger.Info("request successfully handled")
		response.WriteToJson(writer, http.StatusOK, deleteRequest)

	}
}
