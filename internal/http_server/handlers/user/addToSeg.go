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

type Adder interface {
	AddUserToSeg(list []string, id int, ttl time.Time) error
	IfExists(userId int, slug []string) error
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

		if err := adder.IfExists(listOfSegments.Id, listOfSegments.Slug); err != nil {
			var userError storage.UserNotExists
			var segmentError storage.SegmentNotExists
			if errors.As(err, &userError) {
				logger.Error("provided user not exists", userError)
				response.WriteToJson(w, http.StatusConflict, fmt.Sprintf("user %d not exists", userError.Id))
				return
			} else if errors.As(err, &segmentError) {
				logger.Error("provided segment is not active", segmentError)
				response.WriteToJson(w, http.StatusNotFound, fmt.Sprintf("slug %s not exists", segmentError.Slug))
				return
			} else {
				logger.Error("unexpected error from db")
				response.WriteToJson(w, http.StatusInternalServerError, "")
				return
			}
		}

		ttl, err := utils.SetTTL(listOfSegments.TTL)
		if err != nil {
			logger.Error("ttl was not set", ttl)
		}

		if err := adder.AddUserToSeg(listOfSegments.Slug, listOfSegments.Id, ttl); err != nil {
			var segmentErr storage.SegmentAlreadyExistsForUserError
			var segmentExErr storage.SegmentAlreadyExistsError
			if errors.As(err, &segmentErr) {
				logger.Error("segment with this user is already associated", segmentErr)
				response.WriteToJson(w, http.StatusConflict, fmt.Sprintf("segment %s with this user is already associated", segmentErr.Slug))
				return
			} else if errors.As(err, &segmentExErr) {
				logger.Error("segment %s already exists", segmentExErr.Slug)
				response.WriteToJson(w, http.StatusConflict, segmentExErr.Error())
				return
			} else {
				logger.Error("unexpected error from db", err)
				response.WriteToJson(w, http.StatusInternalServerError, "")
				return
			}
		}

		logger.Info("request successfully handled")
		response.WriteToJson(w, http.StatusCreated, listOfSegments)
	}
}
