package report

import (
	"fmt"
	"strings"

	"github.com/xabi93/go-coverage-report/parse"
)

// Markdown is a report generator for markdown.
type Markdown struct{}

// Generate generates from given report a table with coverage data in markdown format.
func (Markdown) Generate(r *parse.Report) string {
	var sb strings.Builder

	// Summary
	sb.WriteString(fmt.Sprintf("Coverage: %.2f%%\n", r.Coverage()))

	// table header
	sb.WriteString(fmt.Sprintf("| %-30s | %-10s | %-10s | %-10s |\n", "File", "Total lines", "Uncovered lines", "Percent"))
	sb.WriteString("| :---- | :----: | :----: | :----: |\n")
	for _, f := range r.Files {
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %.2f%% |\n", f.Name, f.Total, f.Uncovered, f.Coverage()))
	}

	return sb.String()
}
