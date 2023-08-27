package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"

	"avito_service/internal/http_server/handlers/segment/response"
	"avito_service/internal/storage"
	"avito_service/internal/storage/postgres/utils"
	"avito_service/pkg"
)

type Tracker interface {
	GetReport(userId int, startDate, endDate time.Time) ([]storage.Segment, error)
}

func GetCSV(logger *slog.Logger, tracker Tracker) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger = logger.With(
			middleware.GetReqID(request.Context()),
		)

		var CSVLog pkg.CSV

		if err := json.NewDecoder(request.Body).Decode(&CSVLog); err != nil {
			logger.Error("failed to decode request", err)
			response.WriteToJson(writer, http.StatusInternalServerError, err)
			return
		}

		logger.Info("request decoded", slog.Any("request", CSVLog))

		start, end := utils.SetPeriod(CSVLog.Month, CSVLog.Year)

		logger.Info("the period of time set")

		segment, err := tracker.GetReport(CSVLog.UserID, start, end)

		if err != nil {
			var noUser storage.UserNotExists
			var noCSV storage.CSVError
			if errors.As(err, &noUser) {
				logger.Error("user not found", noUser.Id)
				response.WriteToJson(writer, http.StatusNotFound, fmt.Sprintf("user for id: %d not found", noUser.Id))
				return
			} else if errors.As(err, &noCSV) {
				logger.Error("CSV not found", noCSV)
				response.WriteToJson(writer, http.StatusNotFound, noCSV.Error())
				return
			} else {
				logger.Error("unexpected error", err)
				response.WriteToJson(writer, http.StatusInternalServerError, "")
				return
			}
		}

		logger.Info("success", segment)

		response.WriteToJson(writer, http.StatusOK, segment)

	}
}
