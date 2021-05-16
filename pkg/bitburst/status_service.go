package bitburst

import (
	logging2 "bitburst/internal/pkg/logging"
	"bitburst/pkg/online"
	"context"
	"fmt"
	"time"
)

type Service struct {
	online.Client
	Repository online.Upsert
}

func (s Service) Handle(ctx context.Context, ids []int) error {
	status, err := func() ([]online.Status, error) {
		defer logging2.Elapsed(fmt.Sprintf("Client get %d ids", len(ids)))()
		return s.GetStatus(ids)
	}()
	if err != nil {
		return err
	}
	var onlineIds []int
	for _, s := range status {
		if s.Online {
			onlineIds = append(onlineIds, s.Id)
		}
	}
	err = func() error {
		defer logging2.Elapsed(fmt.Sprintf("Repository.UpsertAll %d elements", len(status)))()
		err := s.Repository.UpsertAll(ctx, onlineIds, time.Now())
		return err
	}()
	return err
}
