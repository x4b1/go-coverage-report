package format

import "github.com/xabi93/go-coverage-report/pkg/cover"

type Formatter interface {
	Format(r *cover.Report) (string, error)
}
