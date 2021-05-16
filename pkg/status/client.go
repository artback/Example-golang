package status

import (
	"bitburst/pkg/online"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"sync"
)

type client struct {
	Client  *http.Client
	baseURL string
}

func NewClient(hClient *http.Client, baseURL string) online.Client {
	return &client{Client: hClient, baseURL: baseURL}
}

func (c client) getStatus(id int) (*online.Status, error) {
	url := c.baseURL + "/" + fmt.Sprint(id)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("status %d Get %e", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: status code not OK", id)
	}
	status, err := Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("status %d decode err: %e", id, err)
	}
	return status, nil
}

type result struct {
	status *online.Status
	err    error
}

func unique(arr []int) []int {
	occured := map[int]bool{}
	result := []int{}
	for e := range arr {

		// check if already the mapped
		// variable is set to true or not
		if occured[arr[e]] != true {
			occured[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}
func (c client) GetStatus(ids []int) ([]online.Status, error) {
	rChan := make(chan result, len(ids))
	var wg sync.WaitGroup
	for _, id := range unique(ids) {
		wg.Add(1)
		go func(id int, ch chan result) {
			defer wg.Done()
			status, err := c.getStatus(id)
			ch <- result{status, err}
		}(id, rChan)
	}
	go func() {
		wg.Wait()
		close(rChan)
	}()

	var status []online.Status
	var errs []string
	for r := range rChan {
		if r.err != nil {
			errs = append(errs, r.err.Error())
		}
		if r.status != nil {
			status = append(status, *r.status)
		}
	}
	if errs != nil {
		return status, errors.Wrap(fmt.Errorf(strings.Join(errs, "\n")), "GetStatus:")
	}
	return status, nil
}
