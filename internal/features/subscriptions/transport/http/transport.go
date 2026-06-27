package subshttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/logger"
	httprequest "github.com/shitaiv1ck/test-effective-mobile/internal/core/transport/http/request"
	httpresponse "github.com/shitaiv1ck/test-effective-mobile/internal/core/transport/http/response"
)

type SubsHTTP struct {
	service SubsService
}

type SubsService interface {
	CreateSub(ctx context.Context, sub domains.Sub) (domains.Sub, error)
	GetSubs(ctx context.Context, limit *int, offset *int) ([]domains.Sub, error)
	GetSub(ctx context.Context, subID uuid.UUID) (domains.Sub, error)
	PatchSub(ctx context.Context, subID uuid.UUID, patch domains.PatchSub) (domains.Sub, error)
	DeleteSub(ctx context.Context, subID uuid.UUID) error
}

func NewSubHTTP(service SubsService) *SubsHTTP {
	return &SubsHTTP{
		service: service,
	}
}

func (t *SubsHTTP) CreateSubHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpresponse.NewRW(w)
		logger := logger.FromContext(r.Context())

		logger.Debug("invoke CreateSub handler")

		var request SubDTORequest
		if err := httprequest.DecodeAndValidate(r, &request); err != nil {
			rw.ErrorResponse("failed to decode and validate", err, logger)

			return
		}

		startDate, err := parseStartDate(request.StartDate)
		if err != nil {
			rw.ErrorResponse("failed to parse date", err, logger)

			return
		}

		endDate, err := parseEndDate(request.EndDate)
		if err != nil {
			rw.ErrorResponse("failed to parse date", err, logger)

			return
		}

		sub := domains.Sub{
			ServiceName: request.ServiceName,
			Price:       request.Price,
			UserID:      request.UserID,
			StartDate:   startDate,
			EndDate:     endDate,
		}

		createdSub, err := t.service.CreateSub(r.Context(), sub)
		if err != nil {
			rw.ErrorResponse("failed to create subscription", err, logger)

			return
		}

		response := ToDTO(createdSub)

		rw.JSONResponse(response, http.StatusCreated)
	})
}

func (t *SubsHTTP) GetSubsHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpresponse.NewRW(w)
		logger := logger.FromContext(r.Context())

		logger.Debug("invoke GetSubs handler")

		limit, err := httprequest.GetIntQueryParam(r, "limit")
		if err != nil {
			rw.ErrorResponse("failed to get limit from query param", err, logger)

			return
		}

		offset, err := httprequest.GetIntQueryParam(r, "offset")
		if err != nil {
			rw.ErrorResponse("failed to get offset from query param", err, logger)

			return
		}

		subs, err := t.service.GetSubs(r.Context(), limit, offset)
		if err != nil {
			rw.ErrorResponse("failed to get subscriptions", err, logger)

			return
		}

		response := ToDTOs(subs)

		rw.JSONResponse(response, http.StatusOK)
	})
}

func (t *SubsHTTP) GetSubHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpresponse.NewRW(w)
		logger := logger.FromContext(r.Context())

		logger.Debug("invoke GetSub handler")

		subID, err := httprequest.GetUUIDPathValue(r, "sub_id")
		if err != nil {
			rw.ErrorResponse("failed to get UUID path value", err, logger)

			return
		}

		sub, err := t.service.GetSub(r.Context(), subID)
		if err != nil {
			rw.ErrorResponse("failed to get subscription", err, logger)

			return
		}

		response := ToDTO(sub)

		rw.JSONResponse(response, http.StatusOK)
	})
}

func (t *SubsHTTP) PatchSubHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpresponse.NewRW(w)
		logger := logger.FromContext(r.Context())

		logger.Debug("invoke PatchSub handler")

		subID, err := httprequest.GetUUIDPathValue(r, "sub_id")
		if err != nil {
			rw.ErrorResponse("failed to get UUID path value", err, logger)

			return
		}

		var request PatchSubDTORequest
		if err := httprequest.DecodeAndValidate(r, &request); err != nil {
			rw.ErrorResponse("failed to decode and validate", err, logger)

			return
		}

		patch := domains.PatchSub{Price: request.Price, EndDate: request.EndDate}

		patchedSub, err := t.service.PatchSub(r.Context(), subID, patch)
		if err != nil {
			rw.ErrorResponse("failed to patch subscription", err, logger)

			return
		}

		response := ToDTO(patchedSub)

		rw.JSONResponse(response, http.StatusOK)
	})
}

func (t *SubsHTTP) DeleteSubHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpresponse.NewRW(w)
		logger := logger.FromContext(r.Context())

		logger.Debug("invole DeleteSub handler")

		subID, err := httprequest.GetUUIDPathValue(r, "sub_id")
		if err != nil {
			rw.ErrorResponse("failed to get UUID path value", err, logger)

			return
		}

		if err := t.service.DeleteSub(r.Context(), subID); err != nil {
			rw.ErrorResponse("failed to delete subscription", err, logger)

			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})
}

func parseStartDate(value string) (time.Time, error) {
	date, err := time.Parse("01-2006", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse `%v` to date format `01-2006`: %w", value, errs.ErrInvalidArg)
	}

	return date, err
}

func parseEndDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	date, err := time.Parse("01-2006", *value)
	if err != nil {
		return &time.Time{}, fmt.Errorf("failed to parse `%v` to date format `01-2006`: %w", value, errs.ErrInvalidArg)
	}

	return &date, err
}
