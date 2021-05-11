package online

import (
	"context"
	"time"
)

type Repository interface {
	Upsert
	Delete
}
type Delete interface {
	DeleteOlder(ctx context.Context, time time.Time) error
}
type Upsert interface {
	UpsertAll(ctx context.Context, status []Status, time time.Time) error
}
