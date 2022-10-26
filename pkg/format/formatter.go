package format

import "github.com/x4b1/go-coverage-report/pkg/cover"

type Formatter interface {
	Format(r *cover.Report) (string, error)
}
