package user

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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

		var csvData bytes.Buffer
		csvWriter := csv.NewWriter(&csvData)

		csvWriter.Write([]string{"user_id", "is_active", "created_at", "deleted_at"})

		for _, segment := range segment {
			csvWriter.Write([]string{
				strconv.Itoa(segment.UserId),
				strconv.FormatBool(segment.IsActive),
				segment.CreatedAt.Format("2006-01-02 15:04:05"),
				segment.DeletedAt.Time.Format("2006-01-02 15:04:05"),
			})
		}
		csvWriter.Flush()
		writer.Header().Set("Content-Type", "text/csv")
		writer.Header().Set("Content-Disposition", "attachment; filename=segments_history.csv")

		writer.Write(csvData.Bytes())
	}
}
