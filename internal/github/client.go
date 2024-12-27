package github

import (
	ratelimit "github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v67/github"
	"github.com/karlhepler/disfunction/internal/log"
	"github.com/pkg/errors"
)

type Client struct {
	gh  *github.Client
	log log.Logger
}

func NewClient(ghtoken string, log log.Logger) (*Client, error) {
	gh, err := newGitHubClient(ghtoken)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Client{gh, log}, nil
}

func newGitHubClient(ghtoken string) (*github.Client, error) {
	rateLimiter, err := ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return github.NewClient(rateLimiter).WithAuthToken(ghtoken), nil
}
