package types

import "time"

type Issue struct {
	ID        int       `json:"number"`
	Org       string    `json:"org"`
	Repo      string    `json:"repo"`
	Title     string    `json:"title"`
	URL       string    `json:"html_url"`
	Labels    []Label   `json:"labels"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type Label struct {
	Name string `json:"name"`
}
