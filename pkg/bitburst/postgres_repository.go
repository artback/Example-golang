package bitburst

import (
	"bitburst/pkg/bitburst/migrations"
	"bitburst/pkg/online"
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"strings"
	"time"
)

// Migration version
const version = 1

// Migrate migrates the Postgres schema to the current version.
func ValidateSchema(version uint, db *sql.DB) error {
	sourceInstance, err := bindata.WithInstance(bindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return err
	}
	targetInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("go-bindata", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return err
	}
	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return sourceInstance.Close()
}

const insertStatus = "INSERT INTO status (id,last_seen) VALUES %s ON CONFLICT(id) DO UPDATE SET last_seen = excluded.last_seen"
const deleteLastSeen = `DELETE FROM status WHERE last_seen <= $1`

type sqlContext interface {
	ExecContext(ctx context.Context, valueString string, args ...interface{}) (sql.Result, error)
}
type postgresRepository struct {
	sqlContext
}

func NewPostgresRepository(connString string) (online.Repository, error) {
	c, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parsing postgres URI: %w", err)
	}
	db := stdlib.OpenDB(*c)
	err = ValidateSchema(version, db)
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
	valueString := BuildValuesString(insertStatus, len(ids))
	_, err := p.ExecContext(ctx, valueString, valueArgs...)
	return err
}

// DeleteOlder deletes every status entry older or equal to t(time.time)
func (p postgresRepository) DeleteOlder(ctx context.Context, t time.Time) error {
	_, err := p.ExecContext(ctx, deleteLastSeen, t.Format(time.RFC3339))
	return err
}

func BuildValuesString(strFmt string, length int) string {
	var valueStrings []string
	for i := 0; i < length; i++ {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d)", i*2+1, i*2+2))
	}
	return fmt.Sprintf(strFmt, strings.Join(valueStrings, ","))
}
