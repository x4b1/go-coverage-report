package format

import "github.com/xabi93/go-coverage-report/cover"

type Formatter interface {
	Format(r *cover.Report) (string, error)
}
