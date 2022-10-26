package stdout

import (
	"context"
	"fmt"

	"github.com/xabi93/go-coverage-report/pkg/cover"
)

func NewNotifier() *Notifier {
	return &Notifier{}
}

// Notifier implements coverage report notification for github.
type Notifier struct{}

// Notify creates a check run into github pull request with the given coverage report.
func (n *Notifier) Notify(_ context.Context, _ *cover.Report, body string) error {
	fmt.Println(body)

	return nil
}
