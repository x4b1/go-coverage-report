package parse

import (
	"golang.org/x/tools/cover"
)

func FromLocal(prof []*cover.Profile) (*Report, error) {
	r := Report{
		Files: make([]*File, len(prof)),
	}

	for i, p := range prof {
		f := File{
			Name:  p.FileName,
			Lines: make([]*Line, 0, len(p.Blocks)),
		}

		for _, b := range p.Blocks {
			f.Total += b.NumStmt
			if b.Count < 1 {
				f.Uncovered += b.NumStmt
				f.Lines = append(f.Lines, &Line{
					Start: b.StartLine,
					End:   b.EndLine,
				})
			}
		}

		r.Files[i] = &f

		r.Total += f.Total
		r.Uncovered += f.Uncovered

	}
	return &r, nil
}
