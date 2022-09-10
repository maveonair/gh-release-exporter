package main

import (
	"fmt"
	"github.com/maveonair/gh-release-exporter/internal/metrics"
	"github.com/maveonair/gh-release-exporter/internal/releases"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockGithubClient struct {
}

func (c *MockGithubClient) GetLatestReleases(repo string) ([]string, error) {
	if repo == "redis/redis" {
		return []string{"7.0.4", "7.0.3", "7.0.1"}, nil
	}

	if repo == "strongswan/strongswan" {
		return []string{"5.9.1", "5.9.0"}, nil
	}

	return nil, fmt.Errorf("given repo {%s} is unknown", repo)
}

func Test_CheckReleasesForUpdate(t *testing.T) {
	go metrics.Setup("127.0.0.1:9054")

	mockGithubClient := &MockGithubClient{}

	mockReleases := map[string]releases.Release{
		"redis": {
			LastKnownVersion: "7.0.3",
			GitHubRepo:       "redis/redis",
		},
		"strongswan": {
			LastKnownVersion: "5.9.1",
			GitHubRepo:       "strongswan/strongswan",
		},
	}

	checkReleasesForUpdate(mockGithubClient, mockReleases)

	res, err := http.Get("http://localhost:9054/metrics")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("metrics endpoint could not be fetched (statusCode: %d)", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// Ignore
		}
	}(res.Body)

	data := string(body)

	testSuite := []struct {
		Name          string
		ExpectedValue int
		Payload       string
	}{
		{
			Name:          "Redis",
			ExpectedValue: 0,
			Payload:       "gh_release_probe_success{name=\"redis\"} 0",
		},
		{
			Name:          "Strongswan",
			ExpectedValue: 1,
			Payload:       "gh_release_probe_success{name=\"strongswan\"} 1",
		},
	}

	for _, test := range testSuite {
		if !strings.Contains(data, test.Payload) {
			t.Errorf("Expected probe for %s to be %d", test.Name, test.ExpectedValue)
		}
	}

}
