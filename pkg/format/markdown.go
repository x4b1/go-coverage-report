package format

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/xabi93/go-coverage-report/pkg/cover"
)

//go:embed markdown.tmpl
var defaultTmpl string

var _ Formatter = (*Markdown)(nil)

func NewMarkdown(tmpl string) (*Markdown, error) {
	if tmpl == "" {
		tmpl = defaultTmpl
	}

	t, err := template.New("markdown").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	return &Markdown{t}, nil
}

// Markdown is a report generator for markdown.
type Markdown struct {
	tmpl *template.Template
}

// Generate generates from given report a table with coverage data in markdown format.
func (m *Markdown) Format(r *cover.Report) (string, error) {
	var outBuff bytes.Buffer

	if err := m.tmpl.Execute(&outBuff, r); err != nil {
		return "", err
	}

	return outBuff.String(), nil
}
