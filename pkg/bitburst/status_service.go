package bitburst

import (
	"bitburst/internal/pkg/elapsed"
	"bitburst/pkg/online"
	"context"
	"fmt"
	"time"
)

type Service struct {
	online.Client
	online.Upsert
	elapsed.Log
}

func (s Service) Handle(ctx context.Context, ids []int) error {

	status, err := func() ([]online.Status, error) {
		defer s.Elapsed(fmt.Sprintf("Client get %d ids", len(ids)))()
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
		defer s.Elapsed(fmt.Sprintf("Repository.UpsertAll %d elements", len(status)))()
		err := s.UpsertAll(ctx, onlineIds, time.Now())
		return err
	}()
	return err
}
