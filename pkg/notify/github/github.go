package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/xabi93/go-coverage-report/pkg/cover"
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

func outputTotal(report *cover.Report, gh *githubactions.Action) {
	gh.SetOutput("total", fmt.Sprintf("%.2f", report.Coverage()))
}
