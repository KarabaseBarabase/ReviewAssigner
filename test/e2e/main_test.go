package e2e

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *E2ETestSuite) TestHealthCheck() {
	resp, err := suite.makeRequest("GET", "/health", nil)
	suite.NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var healthResp struct {
		Status  string `json:"status"`
		Service string `json:"service"`
	}
	suite.parseResponse(resp, &healthResp)

	assert.Equal(suite.T(), "ok", healthResp.Status)
	assert.Equal(suite.T(), "review-service", healthResp.Service)
}

func (suite *E2ETestSuite) TestExistingData() {
	resp, err := suite.makeRequest("GET", "/team/get?team_name=backend", nil)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var teamResp struct {
		TeamName string `json:"team_name"`
		Members  []struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			IsActive bool   `json:"is_active"`
		} `json:"members"`
	}
	suite.parseResponse(resp, &teamResp)

	assert.Equal(suite.T(), "backend", teamResp.TeamName)
	assert.Len(suite.T(), teamResp.Members, 3) // u1, u2, u3

	resp, err = suite.makeRequest("GET", "/users/getReview?user_id=u2", nil)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var userPRsResp struct {
		UserID       string `json:"user_id"`
		PullRequests []struct {
			PullRequestID   string `json:"pull_request_id"`
			PullRequestName string `json:"pull_request_name"`
		} `json:"pull_requests"`
	}
	suite.parseResponse(resp, &userPRsResp)

	assert.Equal(suite.T(), "u2", userPRsResp.UserID)
	assert.True(suite.T(), len(userPRsResp.PullRequests) > 0, "User u2 should have assigned PRs")
}

func (suite *E2ETestSuite) TestCreatePRWithExistingAuthor() {
	createPRReq := map[string]interface{}{
		"pull_request_id":   "pr-e2e-new",
		"pull_request_name": "E2E Test New Feature",
		"author_id":         "u1",
	}

	resp, err := suite.makeRequest("POST", "/pullRequest/create", createPRReq)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var prResp struct {
		PR struct {
			PullRequestID     string   `json:"pull_request_id"`
			PullRequestName   string   `json:"pull_request_name"`
			AuthorID          string   `json:"author_id"`
			Status            string   `json:"status"`
			AssignedReviewers []string `json:"assigned_reviewers"`
		} `json:"pr"`
	}
	suite.parseResponse(resp, &prResp)

	assert.Equal(suite.T(), "pr-e2e-new", prResp.PR.PullRequestID)
	assert.Equal(suite.T(), "OPEN", prResp.PR.Status)
	assert.Equal(suite.T(), "u1", prResp.PR.AuthorID)

	assert.Len(suite.T(), prResp.PR.AssignedReviewers, 2)
	for _, reviewer := range prResp.PR.AssignedReviewers {
		assert.NotEqual(suite.T(), "u1", reviewer, "Author should not be self-assigned")
		assert.Contains(suite.T(), []string{"u2", "u3"}, reviewer, "Reviewer should be from backend team")
	}
}

func (suite *E2ETestSuite) TestMergeExistingPR() {
	mergeReq := map[string]interface{}{
		"pull_request_id": "pr-1006",
	}

	resp, err := suite.makeRequest("POST", "/pullRequest/merge", mergeReq)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var mergeResp struct {
		PR struct {
			PullRequestID string `json:"pull_request_id"`
			Status        string `json:"status"`
			MergedAt      string `json:"mergedAt,omitempty"`
		} `json:"pr"`
	}
	suite.parseResponse(resp, &mergeResp)

	assert.Equal(suite.T(), "pr-1006", mergeResp.PR.PullRequestID)
	assert.Equal(suite.T(), "MERGED", mergeResp.PR.Status)
	assert.NotEmpty(suite.T(), mergeResp.PR.MergedAt, "MergedAt should be set")
}

func (suite *E2ETestSuite) TestUserDeactivation() {
	deactivateReq := map[string]interface{}{
		"user_id":   "u3",
		"is_active": false,
	}

	resp, err := suite.makeRequest("POST", "/users/setIsActive", deactivateReq)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var userResp struct {
		User struct {
			UserID   string `json:"user_id"`
			IsActive bool   `json:"is_active"`
		} `json:"user"`
	}
	suite.parseResponse(resp, &userResp)

	assert.Equal(suite.T(), "u3", userResp.User.UserID)
	assert.False(suite.T(), userResp.User.IsActive, "User should be deactivated")
}

func (suite *E2ETestSuite) TestStatistics() {
	resp, err := suite.makeRequest("GET", "/stats/user-assignments", nil)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var statsResp struct {
		UserAssignments map[string]int `json:"user_assignments"`
	}
	suite.parseResponse(resp, &statsResp)

	assert.Greater(suite.T(), statsResp.UserAssignments["u2"], 0, "u2 should have assignments from test data")

	resp, err = suite.makeRequest("GET", "/stats/pr-metrics", nil)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var metricsResp struct {
		PRMetrics map[string]interface{} `json:"pr_metrics"`
	}
	suite.parseResponse(resp, &metricsResp)

	assert.Contains(suite.T(), metricsResp.PRMetrics, "total_prs")
	assert.Contains(suite.T(), metricsResp.PRMetrics, "open_prs")
	assert.Contains(suite.T(), metricsResp.PRMetrics, "merged_prs")

	totalPRs := int(metricsResp.PRMetrics["total_prs"].(float64))
	assert.Greater(suite.T(), totalPRs, 0, "Should have some PRs from test data")
}

func (suite *E2ETestSuite) TestErrorScenarios() {
	invalidPRReq := map[string]interface{}{
		"pull_request_id":   "pr-error-test",
		"pull_request_name": "Error Test",
		"author_id":         "non-existent-user",
	}

	resp, err := suite.makeRequest("POST", "/pullRequest/create", invalidPRReq)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)

	var errorResp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	suite.parseResponse(resp, &errorResp)
	assert.Equal(suite.T(), "NOT_FOUND", errorResp.Error.Code)

	duplicatePRReq := map[string]interface{}{
		"pull_request_id":   "pr-1005",
		"pull_request_name": "Duplicate PR",
		"author_id":         "u1",
	}

	resp, err = suite.makeRequest("POST", "/pullRequest/create", duplicatePRReq)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusConflict, resp.StatusCode)

	suite.parseResponse(resp, &errorResp)
	assert.Equal(suite.T(), "PR_EXISTS", errorResp.Error.Code)
}

func (suite *E2ETestSuite) TestTeamOperations() {
	newTeamReq := map[string]interface{}{
		"team_name": "devops-team",
		"members": []map[string]interface{}{
			{
				"user_id":   "devops-1",
				"username":  "DevOps Engineer 1",
				"is_active": true,
			},
			{
				"user_id":   "devops-2",
				"username":  "DevOps Engineer 2",
				"is_active": true,
			},
		},
	}

	resp, err := suite.makeRequest("POST", "/team/add", newTeamReq)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	resp, err = suite.makeRequest("GET", "/team/get?team_name=devops-team", nil)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var teamResp struct {
		TeamName string `json:"team_name"`
		Members  []struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			IsActive bool   `json:"is_active"`
		} `json:"members"`
	}
	suite.parseResponse(resp, &teamResp)

	assert.Equal(suite.T(), "devops-team", teamResp.TeamName)
	assert.Len(suite.T(), teamResp.Members, 2)
}

func (suite *E2ETestSuite) TestMassDeactivation() {
	massDeactivateReq := map[string]interface{}{
		"user_ids": []string{"u4", "u5"},
	}

	resp, err := suite.makeRequest("POST", "/team/frontend/deactivate-users", massDeactivateReq)
	suite.NoError(err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var massDeactivateResp struct {
		TeamName string                 `json:"team_name"`
		Results  map[string]interface{} `json:"results"`
	}
	suite.parseResponse(resp, &massDeactivateResp)

	assert.Equal(suite.T(), "frontend", massDeactivateResp.TeamName)
	assert.Contains(suite.T(), massDeactivateResp.Results, "u4")
	assert.Contains(suite.T(), massDeactivateResp.Results, "u5")
}

func TestE2ESuite(t *testing.T) {
	if os.Getenv("E2E_TEST") == "" {
		t.Skip("Skipping E2E tests. Set E2E_TEST=1 to run.")
	}
	suite.Run(t, new(E2ETestSuite))
}
