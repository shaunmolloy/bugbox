package types

import "time"

type Issue struct {
	ID			int			`json:"number"`
	Org			string		`json:"org"`
	Title		string		`json:"title"`
	URL			string		`json:"html_url"`
	Labels		[]Label		`json:"labels"`
	CreatedAt	time.Time	`json:"created_at"`
}

type Label struct {
	Name	string		`json:"name"`
}
