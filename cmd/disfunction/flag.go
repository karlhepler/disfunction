package main

import (
	"strings"

	"github.com/karlhepler/disfunction/internal/github"
)

func parseRepos(inputRepos []string) (repos []*github.Repository) {
	if len(inputRepos) == 0 {
		return
	}

	repos = make([]*github.Repository, len(inputRepos))
	for i, inputRepo := range inputRepos {
		repos[i] = &github.Repository{}

		// default to repo name
		if !strings.Contains(inputRepo, "/") {
			inputRepo = "/" + inputRepo
		}

		parts := strings.Split(inputRepo, "/")
		if ownerLogin := parts[0]; ownerLogin != "" {
			repos[i].Owner = &github.User{
				Login: &ownerLogin,
			}
		}
		if len(parts) > 1 {
			if repoName := parts[1]; repoName != "" {
				repos[i].Name = &repoName
			}
		}
	}

	return
}
