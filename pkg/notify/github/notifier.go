package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v48/github"
)

//go:generate moq -stub -pkg github_test -out mocks_test.go . CheckCreator

// CheckCreator is just an interface to allow easy testing
type CheckCreator interface {
	CreateCheckRun(ctx context.Context, owner, repo string, opts github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error)
}

const DefaultCheckName = "Coverage report"

func NewNotifier(cc CheckCreator, owner, repo, headSHA, checkName string) *Notifier {
	if checkName == "" {
		checkName = DefaultCheckName
	}
	return &Notifier{cc, owner, repo, headSHA, checkName}
}

// Notifier implements coverage report notification for github.
type Notifier struct {
	cli     CheckCreator
	owner   string
	repo    string
	headSHA string

	checkName string
}

// Notify creates a check run into github pull request with the given coverage report.
func (n *Notifier) Notify(ctx context.Context, body string) error {
	_, _, err := n.cli.CreateCheckRun(ctx, n.owner, n.repo, github.CreateCheckRunOptions{
		Name:       n.checkName,
		HeadSHA:    n.headSHA,
		Status:     github.String("completed"),
		Conclusion: github.String("neutral"),
		Output: &github.CheckRunOutput{
			Title:   github.String(n.checkName),
			Summary: github.String(body),
		},
	})
	if err != nil {
		return fmt.Errorf("github notify: %w", err)
	}

	return nil
}