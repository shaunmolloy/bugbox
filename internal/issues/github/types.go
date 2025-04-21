package github

import "github.com/shaunmolloy/bugbox/internal/types"

type IssueResponse struct {
	Count	int				`json:"total_count"`
	Items	[]types.Issue	`json:"items"`
}

