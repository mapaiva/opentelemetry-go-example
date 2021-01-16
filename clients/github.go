package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GithubAPI represents a GithubAPI HTTP client.
type GithubAPI struct {
	HTTPClient *http.Client
	URL        string
}

// UserByUsername returns a GH user by username.
func (g GithubAPI) UserByUsername(ctx context.Context, username string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/users/%s", g.URL, username), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var userResponse map[string]interface{}
	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(&userResponse)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return userResponse, nil
}
