package repository

import (
	"bitburst/internal/pkg/postgres"
	"bitburst/pkg/online"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"net/url"
	"time"
)

// Migration version
const version = 1

const insertStatus = "INSERT INTO status (id,last_seen) VALUES %s ON CONFLICT(id) DO UPDATE SET last_seen = excluded.last_seen"
const deleteLastSeen = `DELETE FROM status WHERE last_seen <= $1`

type sqlContext interface {
	ExecContext(ctx context.Context, valueString string, args ...interface{}) (sql.Result, error)
}
type postgresRepository struct {
	sqlContext
}

func NewPostgresRepository(url *url.URL) (online.Repository, error) {
	c, err := pgx.ParseConfig(url.String())
	if err != nil {
		return nil, fmt.Errorf("parsing postgres URI: %w", err)
	}
	db := stdlib.OpenDB(*c)
	err = postgres.ValidateSchema(version, db)
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, db.Ping()
}

// UpsertAll runs a transaction which updates each entry by first deleting and then inserting again
// Uses application time instead of database time for compatibility and testability with DeleteOlder function
func (p postgresRepository) UpsertAll(ctx context.Context, ids []int, t time.Time) error {
	if len(ids) == 0 {
		// Empty input isn't an error but it is unnecessary to continue function
		return nil
	}
	valueArgs := []interface{}{}
	for _, s := range ids {
		valueArgs = append(valueArgs, s)
		valueArgs = append(valueArgs, t.Format(time.RFC3339))
	}
	valueString := postgres.BuildValuesString(insertStatus, len(ids))
	_, err := p.ExecContext(ctx, valueString, valueArgs...)
	return err
}

// DeleteOlder deletes every status entry older or equal to t(time.time)
func (p postgresRepository) DeleteOlder(ctx context.Context, t time.Time) error {
	_, err := p.ExecContext(ctx, deleteLastSeen, t.Format(time.RFC3339))
	return err
}
