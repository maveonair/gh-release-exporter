package github

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type release struct {
	Name string `json:"tag_name"`
}

type Client interface {
	GetLatestReleases(repo string) ([]string, error)
}

type client struct {
	httpClient *http.Client
}

// NewClient returns a new client.
func NewClient() Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetLatestReleases returns the latest releases for the given repo.
func (c *client) GetLatestReleases(repo string) ([]string, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", repo))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tags []release
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, errors.New("no releases found on GitHub")
	}

	versions := make([]string, 0)

	for _, release := range tags {
		versions = append(versions, release.Name)
	}

	return versions, nil
}
