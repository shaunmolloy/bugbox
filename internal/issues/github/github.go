package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/shaunmolloy/bugbox/internal/issues"
	"github.com/shaunmolloy/bugbox/internal/logging"
	"github.com/shaunmolloy/bugbox/internal/storage/config"
	"github.com/shaunmolloy/bugbox/internal/types"
)

const baseURL = "https://api.github.com"

func FetchAllIssues(fetchAll bool, client issues.HttpClient) error {
	conf, _ := config.LoadConfig()

	// Load existing issues
	issuesConf, err := config.LoadIssues()
	if err != nil {
		issuesConf = config.Issues{} // Initialize an empty Issues map
		logging.Error(fmt.Sprintf("Error loading existing issues: %v", err))
	}

	for _, org := range conf.Orgs {
		issues, err := FetchIssues(org, fetchAll, client)
		if err != nil {
			logging.Error(fmt.Sprintf("Error fetching issues for org %s: %v", org, err))
			continue
		}

		// Process and store issues by org/repo/number
		for _, issue := range issues {
			issue.Repo = parseRepo(issue.URL)

			// Ensure maps exist for this org and repo
			if _, ok := issuesConf[issue.Org]; !ok {
				issuesConf[issue.Org] = make(map[string]map[int]types.Issue)
			}
			if _, ok := issuesConf[issue.Org][issue.Repo]; !ok {
				issuesConf[issue.Org][issue.Repo] = make(map[int]types.Issue)
			}

			// Check if issue already exists to preserve Read status
			if existingIssue, exists := issuesConf[issue.Org][issue.Repo][issue.ID]; exists {
				issue.Read = existingIssue.Read
			}

			// Remove issue from conf if state is closed
			if issue.State == types.StateClosed {
				delete(issuesConf[issue.Org][issue.Repo], issue.ID)
				if len(issuesConf[issue.Org][issue.Repo]) == 0 {
					delete(issuesConf[issue.Org], issue.Repo)
				}
			}

			// Store the issue in the hierarchical structure
			issuesConf[issue.Org][issue.Repo][issue.ID] = issue
		}
	}

	if err := config.SaveIssues(issuesConf); err != nil {
		logging.Error(fmt.Sprintf("Error saving issues: %v", err))
		return err
	}

	return nil
}

func FetchIssues(owner string, fetchAll bool, client issues.HttpClient) ([]types.Issue, error) {
	logging.Info(fmt.Sprintf("Searching GitHub issues in org: %s", owner))
	conf, _ := config.LoadConfig()

	query := fmt.Sprintf("org:%s is:issue is:open sort:created-desc", owner)
	encodedQuery := url.QueryEscape(query)

	page := 1
	var allIssues []types.Issue

	for {
		api := fmt.Sprintf("%s/search/issues?q=%s&per_page=100&page=%d", baseURL, encodedQuery, page)
		logging.Debug(fmt.Sprintf("Fetching %s", api))
		req, err := http.NewRequest("GET", api, nil)
		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		req.Header.Set("Authorization", "token "+conf.GitHubToken)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

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

		allIssues = append(allIssues, result.Items...)

		if !fetchAll || len(result.Items) == 0 || page == 10 {
			break
		}

		page++
	}

	logging.Info(fmt.Sprintf("Found %d issues in org: %s", len(allIssues), owner))
	return allIssues, nil
}
