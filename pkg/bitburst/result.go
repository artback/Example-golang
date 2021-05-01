package bitburst

import (
	"bitburst/pkg/online"
	"sync"
)

type result struct {
	status *online.Status
	err    error
}

func getResult(ids []int, client online.Client) chan result {
	rChan := make(chan result)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id int, c chan result) {
			defer wg.Done()
			status, err := client.GetStatus(id)
			c <- result{status, err}
		}(id, rChan)
	}
	go func() {
		wg.Wait()
		close(rChan)
	}()
	return rChan
}

func readStatus(result chan result) ([]online.Status, error) {
	var status []online.Status
	var err error
	for r := range result {
		if r.err != nil {
			err = r.err
		}
		if r.status != nil {
			status = append(status, *r.status)
		}
	}
	if err != nil {
		return nil, err
	}
	return status, nil
}
