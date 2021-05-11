package bitburst

import (
	"bitburst/internal/logging"
	"bitburst/pkg/online"
	"context"
	"fmt"
	"time"
)

type Service struct {
	online.Client
	online.Repository
}

func (s Service) Handle(ctx context.Context, ids []int) error {
	status, err := func() ([]online.Status, error) {
		defer logging.Elapsed(fmt.Sprintf("Client get %d ids", len(ids)))()
		return s.GetStatus(ids)
	}()
	if err != nil {
		return err
	}
	err = func() error {
		defer logging.Elapsed(fmt.Sprintf("Repository.UpsertAll %d elements", len(status)))()
		err := s.UpsertAll(ctx, status, time.Now())
		return err
	}()
	return err
}
