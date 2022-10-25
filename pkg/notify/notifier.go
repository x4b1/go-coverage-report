package notify

import "context"

// Notifier knows how to notify coverage report result
type Notifier interface {
	// Notify sends coverage report to target.
	Notify(ctx context.Context, report string) error
}
