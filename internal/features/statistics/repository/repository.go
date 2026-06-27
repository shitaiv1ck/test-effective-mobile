package statsrep

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/repository/postgres"
)

type StatsRepository struct {
	store postgres.Postgres
}

func NewStatsRepository(store postgres.Postgres) *StatsRepository {
	return &StatsRepository{
		store: store,
	}
}

func (r *StatsRepository) FindSubs(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
	fromDate *time.Time,
	toDate *time.Time,
) ([]domains.Sub, error) {
	ctx, cancel := context.WithTimeout(ctx, r.store.GetTimeout())
	defer cancel()

	var query strings.Builder
	query.WriteString(`
		SELECT * FROM test_task.subscriptions
	`)

	args := make([]any, 0)
	conditions := make([]string, 0)

	if userID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, userID)
	}

	if serviceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", len(args)+1))
		args = append(args, serviceName)
	}

	if fromDate != nil {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", len(args)+1))
		args = append(args, fromDate)
	}

	if toDate != nil {
		conditions = append(conditions, fmt.Sprintf("start_date < $%d", len(args)+1))
		args = append(args, toDate)
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	subs := make([]domains.Sub, 0)
	rows, err := r.store.Query(ctx, query.String(), args...)
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
