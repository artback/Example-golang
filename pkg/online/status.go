package online

import (
	"encoding/json"
	"errors"
	"io"
)

// Use pointers so that DecodeStatus can return error if id or online is not set
type Status struct {
	Id     int  `json:"id"`
	Online bool `json:"online"`
}
type status struct {
	Id     *int  `json:"id"`
	Online *bool `json:"online"`
}

func DecodeStatus(r io.Reader) (*Status, error) {
	var s status
	err := json.NewDecoder(r).Decode(&s)
	if err != nil {
		return nil, err
	}
	if s.Id == nil {
		return nil, errors.New("id is undefined")
	}
	if s.Online == nil {
		return nil, errors.New("online is undefined")
	}

	return &Status{Id: *s.Id, Online: *s.Online}, nil
}
func GetUnique(status []Status, online bool) []Status {
	var unique = make(map[int]bool)
	var st []Status
	for _, s := range status {
		if s.Online == online && !unique[s.Id] {
			unique[s.Id] = true
			st = append(st, s)
		}
	}
	return st
}
