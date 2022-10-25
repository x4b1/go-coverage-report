package github

import (
	"context"

	"github.com/sethvargo/go-githubactions"
)

func NewStepSummary() *StepSummary {
	return &StepSummary{}
}

// StepSummary implements coverage report notification for github actions step summary.
type StepSummary struct{}

// Notify creates a check run into github pull request with the given coverage report.
func (*StepSummary) Notify(_ context.Context, body string) error {
	githubactions.AddStepSummary(body)

	return nil
}
