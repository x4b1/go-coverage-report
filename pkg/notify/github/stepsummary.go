package github

import (
	"context"

	"github.com/sethvargo/go-githubactions"
	"github.com/x4b1/go-coverage-report/pkg/cover"
)

func NewStepSummary(action *githubactions.Action) *StepSummary {
	return &StepSummary{action}
}

// StepSummary implements coverage report notification for github actions step summary.
type StepSummary struct {
	action *githubactions.Action
}

// Notify creates a check run into github pull request with the given coverage report.
func (ss *StepSummary) Notify(_ context.Context, report *cover.Report, body string) error {
	outputTotal(report, ss.action)

	ss.action.AddStepSummary(body)

	return nil
}
