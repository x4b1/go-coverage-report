package parse

type Report struct {
	Total     int
	Uncovered int
	Files     []*File
}

func (r Report) Coverage() float64 {
	return 100 * float64(r.Total-r.Uncovered) / float64(r.Total)
}

type File struct {
	Name      string
	Total     int
	Uncovered int
	Lines     []*Line
}

func (f File) Coverage() float64 {
	return 100 * float64(f.Total-f.Uncovered) / float64(f.Total)
}

type Line struct {
	Start int
	End   int
}
