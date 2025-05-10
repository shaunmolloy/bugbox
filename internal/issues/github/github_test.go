package github

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shaunmolloy/bugbox/internal/issues"
)

func TestFetchAllIssues(t *testing.T) {
	t.Run("returns nil when fetching all issues", func(t *testing.T) {
		fetchAll := true

		client := &issues.ClientMock{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"total_count": 0, "items": []}`)),
				}, nil
			},
		}

		if err := FetchAllIssues(fetchAll, client); err != nil {
			t.Fatal("expected nil, got error")
		}
	})
}

func TestFetchIssues(t *testing.T) {
	t.Run("returns nil when fetching issues", func(t *testing.T) {
		owner := "example"
		fetchAll := false

		client := &issues.ClientMock{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"total_count": 0, "items": []}`)),
				}, nil
			},
		}

		if _, err := FetchIssues(owner, fetchAll, client); err != nil {
			t.Fatal("expected nil, got error")
		}
	})
}
