package github

import "github.com/google/go-github/v67/github"

type User = github.User

type Owner struct {
	Login string
}

func (o Owner) String() string {
	return o.Login
}
