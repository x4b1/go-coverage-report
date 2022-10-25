package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v48/github"
	"github.com/xabi93/go-coverage-report/internal/log"
)

//go:generate moq -stub -pkg github_test -out mocks_test.go . CheckCreator

// CheckCreator is just an interface to allow easy testing
type CheckCreator interface {
	CreateCheckRun(ctx context.Context, owner, repo string, opts github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error)
}

const DefaultCheckName = "Coverage report"

func NewCheckRun(cc CheckCreator, owner, repo, headSHA, checkName string) *CheckRun {
	if checkName == "" {
		checkName = DefaultCheckName
	}
	return &CheckRun{cc, owner, repo, headSHA, checkName}
}

// CheckRun implements coverage report notification for github check runs.
type CheckRun struct {
	cli     CheckCreator
	owner   string
	repo    string
	headSHA string

	checkName string
}

// Notify creates a check run into github pull request with the given coverage report.
func (c *CheckRun) Notify(ctx context.Context, body string) error {
	cr, r, err := c.cli.CreateCheckRun(ctx, c.owner, c.repo, github.CreateCheckRunOptions{
		Name:       c.checkName,
		HeadSHA:    c.headSHA,
		Status:     github.String("completed"),
		Conclusion: github.String("success"),
		Output: &github.CheckRunOutput{
			Title:   github.String(c.checkName),
			Summary: github.String(body),
		},
	})
	if err != nil {
		return fmt.Errorf("github notify: %w", err)
	}

	log.Debugf("check run create response %d", r.StatusCode)
	log.Debugf("check run id %d", cr.GetID())
	log.Debugf("check run url %d", cr.GetURL())

	return nil
}
