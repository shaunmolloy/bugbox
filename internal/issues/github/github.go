package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shaunmolloy/bugbox/internal/logging"
	"github.com/shaunmolloy/bugbox/internal/storage/config"
	"github.com/shaunmolloy/bugbox/internal/types"
)

const baseURL = "https://api.github.com"

func FetchAllIssues() error {
	conf, _ := config.LoadConfig()
	issuesConf := config.Issues{}

	for _, org := range conf.Orgs {
		issues, err := FetchIssues(org)
		if err != nil {
			logging.Error(fmt.Sprintf("Error fetching issues for org %s: %v", org, err))
		}

		issuesConf = append(issuesConf, issues...)

		if err := config.SaveIssues(issuesConf); err != nil {
			logging.Error(fmt.Sprintf("Error saving issues for org %s: %v", org, err))
			return err
		}
	}
	return nil
}

func FetchIssues(owner string) ([]types.Issue, error) {
	logging.Info(fmt.Sprintf("Searching GitHub issues in org: %s", owner))
	conf, _ := config.LoadConfig()

	query := fmt.Sprintf("org:%s is:issue is:open sort:created-desc", owner)
	encodedQuery := url.QueryEscape(query)

	api := fmt.Sprintf("%s/search/issues?q=%s", baseURL, encodedQuery)
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "token "+conf.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var result IssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Set Repo & Org for each issue
	for i := range result.Items {
		result.Items[i].Org = owner
		result.Items[i].Repo = parseRepo(result.Items[i].URL)
		result.Items[i].Read = false
	}

	logging.Info(fmt.Sprintf("Found %d issues in org: %s", result.Count, owner))
	return result.Items, nil
}
