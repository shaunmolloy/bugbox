package types

import (
	"encoding/json"
	"strings"
)

type State int

const (
	StateOpen State = iota
	StateClosed
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	default:
		return "open"
	}
}

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *State) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch strings.ToLower(str) {
	case "closed":
		*s = StateClosed
	default:
		*s = StateOpen
	}
	return nil
}
