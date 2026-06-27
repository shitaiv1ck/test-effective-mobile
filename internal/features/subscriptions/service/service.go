package subssrvc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
)

type SubsService struct {
	rep SubsRespository
}

type SubsRespository interface {
	SaveSub(ctx context.Context, sub domains.Sub) (domains.Sub, error)
	FindSubs(ctx context.Context, limit *int, offset *int) ([]domains.Sub, error)
	FindSubByID(ctx context.Context, subID uuid.UUID) (domains.Sub, error)
	UpdateSub(ctx context.Context, subID uuid.UUID, patch domains.Sub) (domains.Sub, error)
	DeleteSub(ctx context.Context, subID uuid.UUID) error
}

func NewSubsService(rep SubsRespository) *SubsService {
	return &SubsService{
		rep: rep,
	}
}

func (s *SubsService) CreateSub(ctx context.Context, sub domains.Sub) (domains.Sub, error) {
	if err := sub.Validate(); err != nil {
		return domains.Sub{}, fmt.Errorf("failed to validate subscription: %w", err)
	}

	savedSub, err := s.rep.SaveSub(ctx, sub)
	if err != nil {
		return domains.Sub{}, fmt.Errorf("failed to save subscription: %w", err)
	}

	return savedSub, nil
}

func (s *SubsService) GetSubs(ctx context.Context, limit *int, offset *int) ([]domains.Sub, error) {
	if limit != nil && *limit < 0 {
		return []domains.Sub{}, fmt.Errorf("limit must be non negative: %w", errs.ErrInvalidArg)
	}

	if offset != nil && *limit < 0 {
		return []domains.Sub{}, fmt.Errorf("offset must be non negative: %w", errs.ErrInvalidArg)
	}

	subs, err := s.rep.FindSubs(ctx, limit, offset)
	if err != nil {
		return []domains.Sub{}, fmt.Errorf("failed to find subscriptions: %w", err)
	}

	return subs, err
}

func (s *SubsService) GetSub(ctx context.Context, subID uuid.UUID) (domains.Sub, error) {
	if subID == uuid.Nil {
		return domains.Sub{}, fmt.Errorf("empty sub ID: %w", errs.ErrInvalidArg)
	}

	sub, err := s.rep.FindSubByID(ctx, subID)
	if err != nil {
		return domains.Sub{}, fmt.Errorf("failed to find subscription: %w", err)
	}

	return sub, nil
}

func (s *SubsService) PatchSub(ctx context.Context, subID uuid.UUID, patch domains.PatchSub) (domains.Sub, error) {
	if subID == uuid.Nil {
		return domains.Sub{}, fmt.Errorf("empty sub ID: %w", errs.ErrInvalidArg)
	}

	sub, err := s.rep.FindSubByID(ctx, subID)
	if err != nil {
		return domains.Sub{}, fmt.Errorf("failed to find subscription: %w", err)
	}

	if err := sub.ApplyPatch(patch); err != nil {
		return domains.Sub{}, fmt.Errorf("failed to apply patch to subscription: %w", err)
	}

	updatedSub, err := s.rep.UpdateSub(ctx, subID, sub)
	if err != nil {
		return domains.Sub{}, fmt.Errorf("failed to update subscription: %w", err)
	}

	return updatedSub, nil
}

func (s *SubsService) DeleteSub(ctx context.Context, subID uuid.UUID) error {
	if subID == uuid.Nil {
		return fmt.Errorf("empty sub ID: %w", errs.ErrInvalidArg)
	}

	if err := s.rep.DeleteSub(ctx, subID); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
