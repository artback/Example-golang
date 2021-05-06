package online

import (
	"encoding/json"
	"errors"
	"io"
)

// Use pointers so that DecodeStatus can return error if id or online is not set
type Status struct {
	Id     *int  `json:"id"`
	Online *bool `json:"online"`
}

func NewStatus(id int, online bool) *Status {
	return &Status{
		Id:     &id,
		Online: &online,
	}
}
func DecodeStatus(r io.Reader) (*Status, error) {
	var status Status
	err := json.NewDecoder(r).Decode(&status)
	if err != nil {
		return nil, err
	}
	if status.Id == nil {
		return nil, errors.New("id is undefined")
	}
	if status.Online == nil {
		return nil, errors.New("online is undefined")
	}
	return &status, nil
}
func GetUnique(status []Status, online bool) []Status {
	var unique = make(map[int]bool)
	var st []Status
	for _, s := range status {
		if *s.Online == online && !unique[*s.Id] {
			unique[*s.Id] = true
			st = append(st, s)
		}
	}
	return st
}
