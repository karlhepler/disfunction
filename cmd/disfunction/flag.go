package main

import (
	"github.com/karlhepler/disfunction/internal/github"
)

func parseRepos(inputRepos []string) (repos []*github.Repository) {
	if len(inputRepos) == 0 {
		return
	}

	repos = make([]*github.Repository, len(inputRepos))
	for i, inputRepo := range inputRepos {
		repos[i] = github.NewRepository(inputRepo)
	}

	return
}
