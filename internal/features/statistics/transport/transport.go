package statshttp

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/logger"
	httprequest "github.com/shitaiv1ck/test-effective-mobile/internal/core/transport/http/request"
	httpresponse "github.com/shitaiv1ck/test-effective-mobile/internal/core/transport/http/response"
)

type StatsHTTP struct {
	service StatsService
}

type StatsService interface {
	GetStatistics(
		ctx context.Context,
		userID *uuid.UUID,
		serviceName *string,
		fromDate *time.Time,
		toDate *time.Time,
	) (domains.Statistics, error)
}

func NewStatsHTTP(service StatsService) *StatsHTTP {
	return &StatsHTTP{
		service: service,
	}
}

func (t *StatsHTTP) GetStatisticsHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpresponse.NewRW(w)
		logger := logger.FromContext(r.Context())

		logger.Debug("invoke GetStatistics handler")

		userID, err := httprequest.GetUUIDQueryParam(r, "user_id")
		if err != nil {
			rw.ErrorResponse("failed to get `user_id` from query param", err, logger)

			return
		}

		serviceName := httprequest.GetStringQueryParam(r, "service_name")

		fromDate, err := httprequest.GetDateQueryParam(r, "from_date")
		if err != nil {
			rw.ErrorResponse("failed to get `from_date` from query param", err, logger)

			return
		}

		toDate, err := httprequest.GetDateQueryParam(r, "to_date")
		if err != nil {
			rw.ErrorResponse("failed to get `to_date` from query param", err, logger)

			return
		}

		stats, err := t.service.GetStatistics(r.Context(), userID, serviceName, fromDate, toDate)
		if err != nil {
			rw.ErrorResponse("failed to get statistics", err, logger)

			return
		}

		response := ToDTO(stats)

		rw.JSONResponse(response, http.StatusOK)
	})
}
