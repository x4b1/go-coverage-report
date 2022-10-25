package github

import (
	"context"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func NewClient(ctx context.Context, token string) *github.Client {
	return github.NewClient(
		oauth2.NewClient(
			ctx,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		),
	)
}
