package statssrvc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
)

type StatsService struct {
	rep StatsRepository
}

type StatsRepository interface {
	FindSubs(
		ctx context.Context,
		userID *uuid.UUID,
		serviceName *string,
		fromDate *time.Time,
		toDate *time.Time,
	) ([]domains.Sub, error)
}

func NewStatsService(rep StatsRepository) *StatsService {
	return &StatsService{
		rep: rep,
	}
}

func (s *StatsService) GetStatistics(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	fromDate *time.Time,
	toDate *time.Time,
) (domains.Statistics, error) {
	if userID != nil && *userID == uuid.Nil {
		return domains.Statistics{}, fmt.Errorf("empty sub ID: %w", errs.ErrInvalidArg)
	}

	if serviceName != nil && strings.TrimSpace(*serviceName) == "" {
		return domains.Statistics{}, fmt.Errorf("empty service name: %w", errs.ErrInvalidArg)
	}

	if fromDate != nil && toDate != nil && !fromDate.Before(*toDate) {
		return domains.Statistics{}, fmt.Errorf("`from_date` must be before `to_date`: %w", errs.ErrInvalidArg)
	}

	subs, err := s.rep.FindSubs(ctx, userID, serviceName, fromDate, toDate)
	if err != nil {
		return domains.Statistics{}, fmt.Errorf("failed to find subs: %w", err)
	}

	var stats domains.Statistics
	stats.CalcStatistics(subs)

	return stats, nil
}
