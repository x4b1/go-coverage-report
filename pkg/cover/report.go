package cover

type Report struct {
	TotalStmts int
	Uncovered  int
	Files      []*File
}

func (r Report) Coverage() float64 {
	return 100 * float64(r.TotalStmts-r.Uncovered) / float64(r.TotalStmts)
}

type File struct {
	Name       string
	TotalStmts int
	Uncovered  int
	Lines      []*Line
}

func (f File) Coverage() float64 {
	return 100 * float64(f.TotalStmts-f.Uncovered) / float64(f.TotalStmts)
}

type Line struct {
	Start int
	End   int
}
