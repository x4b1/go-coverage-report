package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/xabi93/go-coverage-report/internal/log"
	"github.com/xabi93/go-coverage-report/pkg/cover"
)

var ErrMissingWorkflowRunField = errors.New("event of type 'workflow_run' is missing 'workflow_run' field")

//go:generate moq -stub -pkg github_test -out mocks_test.go . CheckCreator

// CheckCreator is just an interface to allow easy testing
type CheckCreator interface {
	CreateCheckRun(ctx context.Context, owner, repo string, opts github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error)
}

const DefaultCheckRunName = "Coverage report"

func NewCheckRun(ghAction *githubactions.Action, cc CheckCreator, checkName string) *CheckRun {
	if checkName == "" {
		checkName = DefaultCheckRunName
	}
	return &CheckRun{cc, ghAction, checkName}
}

// CheckRun implements coverage report notification for github check runs.
type CheckRun struct {
	cc CheckCreator

	ghAction *githubactions.Action

	checkName string
}

// Notify creates a check run into github pull request with the given coverage report.
func (c *CheckRun) Notify(ctx context.Context, report *cover.Report, body string) error {
	ghCtx, err := c.ghAction.Context()
	if err != nil {
		return err
	}

	sha, err := c.getSHA(ghCtx)
	if err != nil {
		return err
	}

	owner, repo := ghCtx.Repo()

	outputTotal(report, c.ghAction)

	cr, r, err := c.cc.CreateCheckRun(ctx, owner, repo, github.CreateCheckRunOptions{
		Name:       c.checkName,
		HeadSHA:    sha,
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

	if r != nil {
		log.Debugf("check run create response %d", r.StatusCode)
	}
	if cr != nil {
		log.Debugf("check run id %d", cr.GetID())
		log.Debugf("check run url %s", cr.GetURL())
	}

	return nil
}

func (c *CheckRun) getSHA(ghCtx *githubactions.GitHubContext) (string, error) {
	_, hasPRPayload := ghCtx.Event["pull_request"]
	switch {
	case ghCtx.EventName == "workflow_run":
		c.ghAction.Debugf("Action was triggered by workflow_run: using SHA and RUN_ID from triggering workflow")
		evPayload, ok := ghCtx.Event["workflow_run"]
		if !ok {
			return "", ErrMissingWorkflowRunField
		}

		hc := evPayload.(map[string]any)["head_commit"].(map[string]any)

		return hc["id"].(string), nil

	case hasPRPayload:
		c.ghAction.Debugf("Action was triggered by %s: using SHA from head of source branch", ghCtx.EventName)

		return ghCtx.Event["pull_request"].(map[string]any)["head"].(map[string]any)["sha"].(string), nil
	}

	return ghCtx.SHA, nil
}
