package subsrep

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
	errs "github.com/shitaiv1ck/test-effective-mobile/internal/core/errors"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/repository/postgres"
)

type SubsRepository struct {
	store postgres.Postgres
}

func NewSubsRepository(store postgres.Postgres) *SubsRepository {
	return &SubsRepository{
		store: store,
	}
}

func (r *SubsRepository) SaveSub(ctx context.Context, sub domains.Sub) (domains.Sub, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		INSERT INTO test_task.subscriptions(service_name, price, user_id, start_date, end_date)
		VALUES($1, $2, $3, $4, $5)
		RETURNING *;
	`

	var savedSubs domains.Sub
	if err := r.store.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(
		&savedSubs.ID,
		&savedSubs.ServiceName,
		&savedSubs.Price,
		&savedSubs.UserID,
		&savedSubs.StartDate,
		&savedSubs.EndDate,
	); err != nil {
		return domains.Sub{}, err
	}

	return savedSubs, nil
}

func (r *SubsRepository) FindSubs(ctx context.Context, limit *int, offset *int) ([]domains.Sub, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT * FROM test_task.subscriptions
		LIMIT $1
		OFFSET $2;
	`

	subs := make([]domains.Sub, 0)
	rows, err := r.store.Query(ctx, query, limit, offset)
	if err != nil {
		return []domains.Sub{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub domains.Sub
		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return []domains.Sub{}, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *SubsRepository) FindSubByID(ctx context.Context, subID uuid.UUID) (domains.Sub, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		SELECT * FROM test_task.subscriptions
		WHERE id = $1;
	`

	var sub domains.Sub
	if err := r.store.QueryRow(
		ctx,
		query,
		subID,
	).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domains.Sub{}, fmt.Errorf("subscription with id=%v doesn't exist: %w", subID, errs.ErrNotFound)
		}

		return domains.Sub{}, err
	}

	return sub, nil
}

func (r *SubsRepository) UpdateSub(ctx context.Context, subID uuid.UUID, patch domains.Sub) (domains.Sub, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		UPDATE test_task.subscriptions
		SET price = $1, end_date = $2
		WHERE id = $3
		RETURNING *;
	`

	var updatedSub domains.Sub
	if err := r.store.QueryRow(
		ctx,
		query,
		patch.Price,
		patch.EndDate,
		subID,
	).Scan(
		&updatedSub.ID,
		&updatedSub.ServiceName,
		&updatedSub.Price,
		&updatedSub.UserID,
		&updatedSub.StartDate,
		&updatedSub.EndDate,
	); err != nil {
		return domains.Sub{}, err
	}

	return updatedSub, nil
}

func (r *SubsRepository) DeleteSub(ctx context.Context, subID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	query := `
		DELETE FROM test_task.subscriptions
		WHERE id = $1;
	`

	result, err := r.store.Exec(ctx, query, subID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("subscription with id=%v doesn't exist: %w", subID, errs.ErrNotFound)
	}

	return nil
}
