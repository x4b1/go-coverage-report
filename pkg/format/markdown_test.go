package format_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xabi93/go-coverage-report/pkg/cover"
	"github.com/xabi93/go-coverage-report/pkg/fixtures"
	"github.com/xabi93/go-coverage-report/pkg/format"
)

var testCoverage = &cover.Report{
	TotalStmts: 3,
	Uncovered:  1,
	Files: []*cover.File{
		{
			Name:       "github.com/xabi93/go-coverage-report/fixtures/even.go",
			TotalStmts: 3,
			Uncovered:  1,
			Lines: []*cover.Line{
				{
					Start: 10,
					End:   12,
				},
			},
		},
	},
}

func TestMarkdown(t *testing.T) {
	for _, tc := range []struct {
		name           string
		expectedOutput string
		tmpl           string
	}{
		{
			name:           "default template",
			tmpl:           "",
			expectedOutput: fixtures.DefaultMDResult,
		},
		{
			name:           "custom template",
			tmpl:           `Coverage: {{ printf "%.2f" .Coverage }}%`,
			expectedOutput: "Coverage: 66.67%",
		},
	} {
		f, err := format.NewMarkdown(tc.tmpl)
		require.NoError(t, err)

		out, err := f.Format(testCoverage)
		require.NoError(t, err)

		require.Equal(t, tc.expectedOutput, out)
	}
}
