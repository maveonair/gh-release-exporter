package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type release struct {
	Name string `json:"tag_name"`
}

// GetLatestReleases returns the latest releases for the given repo.
func GetLatestReleases(repo string) ([]string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", repo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
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
