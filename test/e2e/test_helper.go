package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	baseURL string
	client  *http.Client
}

func (suite *E2ETestSuite) SetupSuite() {
	suite.baseURL = "http://localhost:8080"
	suite.client = &http.Client{
		Timeout: 30 * time.Second,
	}

	suite.waitForService()
}

func (suite *E2ETestSuite) waitForService() {
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(suite.baseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			fmt.Println("Service is ready")
			return
		}
		fmt.Printf("⏳ Waiting for service... (attempt %d/%d)\n", i+1, maxAttempts)
		time.Sleep(2 * time.Second)
	}
	suite.FailNow("Service did not become ready in time")
}

func (suite *E2ETestSuite) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		suite.NoError(err)
	}

	req, err := http.NewRequest(method, suite.baseURL+path, bytes.NewBuffer(reqBody))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")

	return suite.client.Do(req)
}

func (suite *E2ETestSuite) parseResponse(resp *http.Response, target interface{}) {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(target)
	suite.NoError(err)
}
