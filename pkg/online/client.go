package online

import (
	"fmt"
	"net/http"
)

type Client interface {
	GetStatus(id int) (*Status, error)
}
type client struct {
	Client  *http.Client
	baseURL string
}

func NewClient(hClient *http.Client, baseURL string) Client {
	return &client{Client: hClient, baseURL: baseURL}
}

func (c client) GetStatus(id int) (*Status, error) {
	url := c.baseURL + "/" + fmt.Sprint(id)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("status %d Get %e", id, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: status code not OK", id)
	}
	status, err := DecodeStatus(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("status %d decode err: %e", id, err)
	}
	return status, nil
}
