package bitburst

import (
	"bitburst/internal/db"
	"bitburst/pkg/online"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/url"
	"time"
)

const schema = `
CREATE TABLE IF NOT EXISTS status(
	id INTEGER PRIMARY KEY,
    last_seen timestamp NOT NULL DEFAULT NOW()
);`
const insertStatus = "INSERT INTO status (id,last_seen) VALUES %s ON CONFLICT(id) DO UPDATE SET last_seen = excluded.last_seen"
const deleteLastSeen = `DELETE FROM status WHERE last_seen <= $1`

type postgresRepository struct {
	*sql.DB
}

func NewPostgresRepository(url *url.URL) (online.Repository, error) {
	db, err := sql.Open("pgx", url.String())
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, err
}

// UpsertAll runs a transaction which updates each entry by first deleting and then inserting again
// Uses application time instead of database time for compatibility and testability with DeleteOlder function
func (p postgresRepository) UpsertAll(ctx context.Context, status []online.Status, t time.Time) error {
	if len(status) == 0 {
		// Empty input isn't an error but it is unnecessary to continue function
		return nil
	}
	valueArgs := []interface{}{}
	for _, s := range status {
		if s.Online {
			valueArgs = append(valueArgs, s.Id)
			valueArgs = append(valueArgs, t.Format(time.RFC3339))
		}
	}
	valueString := db.BuildValuesString(insertStatus, len(status))
	_, err := p.ExecContext(ctx, valueString, valueArgs...)
	return err
}

// DeleteOlder deletes every status entry older or equal to t(time.time)
func (p postgresRepository) DeleteOlder(ctx context.Context, t time.Time) error {
	_, err := p.ExecContext(ctx, deleteLastSeen, t.Format(time.RFC3339))
	return err
}
