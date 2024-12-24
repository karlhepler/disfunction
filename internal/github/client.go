package github

import (
	ratelimit "github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v67/github"
	"github.com/pkg/errors"
)

type Client struct {
	GitHub *github.Client
}

func NewClient(ghtoken string) (*Client, error) {
	gh, err := newGitHubClient(ghtoken)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client := &Client{
		GitHub: gh,
	}

	return client, nil
}

func newGitHubClient(ghtoken string) (*github.Client, error) {
	rl, err := ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return github.NewClient(rl).WithAuthToken(ghtoken), nil
}
