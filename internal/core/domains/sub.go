package domains

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
)

type Sub struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

func (s *Sub) Validate() error {
	if len([]rune(s.ServiceName)) == 0 {
		return fmt.Errorf("len service name must be greater than 0: %w", errs.ErrInvalidArg)
	}

	if s.Price <= 0 {
		return fmt.Errorf("price must be non negative: %w", errs.ErrInvalidArg)
	}

	if s.EndDate != nil && !s.StartDate.Before(*s.EndDate) {
		return fmt.Errorf("start date must be before end date: %w", errs.ErrInvalidArg)
	}

	return nil
}

func (s *Sub) ApplyPatch(patch PatchSub) error {
	if err := patch.Validate(); err != nil {
		return err
	}

	tmp := *s

	if patch.Price.Set && patch.Price.Value != nil {
		tmp.Price = *patch.Price.Value
	}

	if patch.EndDate.Set {
		if patch.EndDate.Value != nil {
			date, _ := time.Parse("01-2006", *patch.EndDate.Value)
			if !date.After(tmp.StartDate) {
				return fmt.Errorf("end date must be after start date: %w", errs.ErrInvalidArg)
			}

			tmp.EndDate = &date
		} else {
			tmp.EndDate = nil
		}
	}

	*s = tmp

	return nil
}

type PatchSub struct {
	Price   Nullable[int]
	EndDate Nullable[string]
}

func (p *PatchSub) Validate() error {
	if p.Price.Value != nil && *p.Price.Value < 0 {
		return fmt.Errorf("price must be non negative: %w", errs.ErrInvalidArg)
	}

	if p.EndDate.Value != nil {
		_, err := time.Parse("01-2006", *p.EndDate.Value)
		if err != nil {
			return fmt.Errorf("failed to parse `%v` to date format `01-2006`: %w", *p.EndDate.Value, errs.ErrInvalidArg)
		}
	}

	return nil
}
