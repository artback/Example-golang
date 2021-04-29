package online

import (
	"context"
	"time"
)

type Repository interface {
	UpsertAll(ctx context.Context, status []Status) error
	DeleteOlder(ctx context.Context, time time.Time) error
}
