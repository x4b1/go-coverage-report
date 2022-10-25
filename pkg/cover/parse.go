package cover

import (
	"context"

	"golang.org/x/tools/cover"
)

type Parser interface {
	Parse(context.Context) (*Report, error)
}

var _ Parser = (*FileParser)(nil)

func NewFileParser(filePath string) *FileParser {
	return &FileParser{filePath}
}

type FileParser struct {
	filePath string
}

func (fp *FileParser) Parse(_ context.Context) (*Report, error) {
	prof, err := cover.ParseProfiles(fp.filePath)
	if err != nil {
		return nil, err
	}

	return parse(prof), nil
}

func parse(prof []*cover.Profile) *Report {
	r := Report{
		Files: make([]*File, len(prof)),
	}

	for i, p := range prof {
		f := File{
			Name:  p.FileName,
			Lines: make([]*Line, 0, len(p.Blocks)),
		}

		for _, b := range p.Blocks {
			f.TotalStmts += b.NumStmt
			if b.Count < 1 {
				f.Uncovered += b.NumStmt
				f.Lines = append(f.Lines, &Line{
					Start: b.StartLine,
					End:   b.EndLine,
				})
			}
		}

		r.Files[i] = &f

		r.TotalStmts += f.TotalStmts
		r.Uncovered += f.Uncovered

	}
	return &r
}
