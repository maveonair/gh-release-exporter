package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type tag struct {
	Name string `json:"name"`
}

// GetLatestTags returns the latest tags for the given repo.
func GetLatestTags(repo string) ([]string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("https://api.github.com/repos/%s/tags", repo))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tags []tag
	err = json.Unmarshal(body, &tags)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, errors.New("No tags found on GitHub")
	}

	versions := make([]string, 0)

	for _, release := range tags {
		versions = append(versions, release.Name)
	}

	return versions, nil
}
