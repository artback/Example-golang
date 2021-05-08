package status

import (
	"bitburst/pkg/online"
	"encoding/json"
	"errors"
	"io"
)

type status struct {
	Id     *int  `json:"id"`
	Online *bool `json:"online"`
}

func Decode(r io.Reader) (*online.Status, error) {
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

	return &online.Status{Id: *s.Id, Online: *s.Online}, nil
}
