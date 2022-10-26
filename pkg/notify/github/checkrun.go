package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golangci/revgrep"
	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-githubactions"
	"github.com/x4b1/go-coverage-report/pkg/cover"
)

var ErrMissingWorkflowRunField = errors.New("event of type 'workflow_run' is missing 'workflow_run' field")

//go:generate moq -stub -pkg github_test -out mocks_test.go . CheckCreator

// Interfaces to allow easy tests github sdk calls.
type (
	CheckCreator interface {
		CreateCheckRun(ctx context.Context, owner, repo string, opts github.CreateCheckRunOptions) (*github.CheckRun, *github.Response, error)
	}
	PRGetter interface {
		GetRaw(ctx context.Context, owner string, repo string, number int, opts github.RawOptions) (string, *github.Response, error)
	}
)

type eventPayload map[string]any

func (ep eventPayload) field(key string) eventPayload {
	return ep[key].(map[string]any)
}

func (ep eventPayload) string(key string) string {
	return ep[key].(string)
}

func (ep eventPayload) int(key string) int {
	return int(ep[key].(float64))
}

const DefaultCheckRunName = "Coverage report"

func NewCheckRun(ghAction *githubactions.Action, cc CheckCreator, pr PRGetter, checkName string) *CheckRun {
	if checkName == "" {
		checkName = DefaultCheckRunName
	}
	return &CheckRun{cc, pr, ghAction, checkName}
}

// CheckRun implements coverage report notification for github check runs.
type CheckRun struct {
	cc CheckCreator
	pr PRGetter

	ghAction *githubactions.Action

	checkName string
}

// Notify creates a check run into github pull request with the given coverage report.
func (c *CheckRun) Notify(ctx context.Context, report *cover.Report, body string) error {
	//nolint:contextcheck //Is not golang context
	ghCtx, err := c.ghAction.Context()
	if err != nil {
		return err
	}

	owner, repo := ghCtx.Repo()
	req := github.CreateCheckRunOptions{
		Name:       c.checkName,
		Status:     github.String("completed"),
		Conclusion: github.String("success"),
		Output: &github.CheckRunOutput{
			Title:   github.String(c.checkName),
			Summary: github.String(body),
		},
	}
	fmt.Println(ghCtx.EventName)
	switch ghCtx.EventName {
	case "pull_request":
		c.ghAction.Debugf("Action was triggered by %s: using SHA from head of source branch", ghCtx.EventName)
		if err = c.handlePR(ctx, ghCtx, &req, report); err != nil {
			return err
		}
	case "workflow_run":
		c.ghAction.Debugf("Action was triggered by workflow_run: using SHA and RUN_ID from triggering workflow")
		req.HeadSHA = eventPayload(ghCtx.Event).field("workflow_run").field("head").string("sha")
	default:
		req.HeadSHA = ghCtx.SHA
	}

	outputTotal(report, c.ghAction)

	_, _, err = c.cc.CreateCheckRun(ctx, owner, repo, req)
	if err != nil {
		return fmt.Errorf("github notify: %w", err)
	}

	return nil
}

func (c *CheckRun) handlePR(ctx context.Context, ghCtx *githubactions.GitHubContext, req *github.CreateCheckRunOptions, report *cover.Report) error {
	payload := eventPayload(ghCtx.Event).field("pull_request")

	req.HeadSHA = payload.field("head").string("sha")

	owner, repo := ghCtx.Repo()

	patch, _, err := c.pr.GetRaw(ctx, owner, repo, payload.int("number"), github.RawOptions{Type: github.Diff})
	if err != nil {
		return err
	}

	checker := revgrep.Checker{Patch: strings.NewReader(patch)}

	if err := checker.Prepare(); err != nil {
		return fmt.Errorf("can't prepare diff by revgrep: %s", err)
	}

	for _, f := range report.Files {
		for _, l := range f.Lines {
			filePath := strings.TrimPrefix(f.Name, "github.com/x4b1/go-coverage-report/")
			_, isNew := checker.IsNewIssue(checkInput{filePath, l.Start})
			if isNew {
				req.Output.Annotations = append(req.Output.Annotations, &github.CheckRunAnnotation{
					Message:         github.String("Uncovered lines"),
					AnnotationLevel: github.String("failure"),
					Path:            github.String(filePath),
					StartLine:       github.Int(l.Start),
					EndLine:         github.Int(l.End),
				})
			}
		}
	}

	return nil
}

type checkInput struct {
	path string
	line int
}

func (ci checkInput) FilePath() string {
	return ci.path
}

func (ci checkInput) Line() int {
	return ci.line
}
