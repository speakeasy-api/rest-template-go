package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
	"github.com/speakeasy-api/rest-template-go/internal/users/model"
	"go.uber.org/zap"
)

// FindUsers will retrieve a list of users based on matching all of the the provided filters and using pagination if limit is gt 0
// Note: depending on the actual use cases for such functionality I would probably take the route of using elasticsearch and opening up
// the flexibility of having a search type function.
func (s *Store) FindUsers(ctx context.Context, filters []model.Filter, offset, limit int64) ([]*model.User, error) {
	if len(filters) == 0 {
		return nil, ErrInvalidFilters.Wrap(errors.ErrInvalidRequest)
	}

	whereClauses := []string{}
	values := []interface{}{}

	for i, f := range filters {
		whereClauses = append(whereClauses, fmt.Sprintf("%s %s $%d", f.Field, f.MatchType, i+1))
		values = append(values, getFindValue(f))
	}

	limitClause := ""

	if limit > 0 {
		limitClause = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := s.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM users WHERE %s ORDER BY id ASC%s", strings.Join(whereClauses, " AND "), limitClause), values...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound.Wrap(err)
		}

		return nil, errors.ErrUnknown.Wrap(err)
	}
	if rows == nil {
		return nil, errors.ErrUnknown
	}
	defer rows.Close()

	users := []*model.User{}

	for rows.Next() {
		var u model.User
		if err := rows.StructScan(&u); err != nil {
			logging.From(ctx).Error("failed to deserialize user from database", zap.Error(err))
		} else {
			users = append(users, &u)
		}
	}

	if len(users) == 0 {
		return nil, errors.ErrNotFound
	}

	return users, nil
}

func getFindValue(f model.Filter) string {
	switch f.MatchType {
	case model.MatchTypeLike:
		return fmt.Sprintf("%%%s%%", f.Value)
	case model.MatchTypeEqual:
		fallthrough
	default:
		return f.Value
	}
}
