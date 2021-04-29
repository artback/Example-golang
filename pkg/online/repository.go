package online

import (
	"context"
	"time"
)

type Repository interface {
	UpsertAll(ctx context.Context, status []Status, time time.Time) error
	DeleteOlder(ctx context.Context, time time.Time) error
}
