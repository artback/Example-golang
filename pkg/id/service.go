package id

import (
	"context"
)

type Service interface {
	Handle(ctx context.Context, ids []int) error
}
