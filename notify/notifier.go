package notify

import "context"

type Notifier interface {
	Notify(ctx context.Context, report string) error
}
