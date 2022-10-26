package notify

import (
	"context"

	"github.com/x4b1/go-coverage-report/pkg/cover"
)

// Notifier knows how to notify coverage report result.
type Notifier interface {
	// Notify sends coverage report to target.
	Notify(ctx context.Context, report *cover.Report, body string) error
}
