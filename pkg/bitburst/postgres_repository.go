package bitburst

import (
	"bitburst/pkg/online"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/url"
	"time"
)

const schema = `
CREATE TABLE IF NOT EXISTS status(
	id text PRIMARY KEY,
    last_seen timestamp NOT NULL DEFAULT NOW()
);`
const insertStatus = "INSERT INTO status (id,last_seen) VALUES ($1,$2)"
const deleteStatus = "DELETE from status where id = $1 "

type postgresRepository struct {
	*sql.DB
}

func NewPostgresRepository(url *url.URL) (online.Repository, error) {
	db, err := sql.Open("postgres", url.String())
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
		return fmt.Errorf("status is empty")
	}
	tx, err := p.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	dStmt, err := tx.PrepareContext(ctx, deleteStatus)
	if err != nil {
		return err
	}
	iStmt, err := tx.PrepareContext(ctx, insertStatus)
	if err != nil {
		return err
	}
	for _, s := range status {
		var err error
		if *s.Online == true {
			_, err = dStmt.ExecContext(ctx, s.Id)
			if err != nil {
				return err
			}
			_, err = iStmt.ExecContext(ctx, s.Id, t.Format(time.RFC3339))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteOlder deletes every status entry older or equal to t(time.time)
func (p postgresRepository) DeleteOlder(ctx context.Context, t time.Time) error {
	_, err := p.ExecContext(ctx, `DELETE FROM status WHERE last_seen <= $1`, t.Format(time.RFC3339))
	return err
}
