package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xabi93/go-coverage-report/input"
	"github.com/xabi93/go-coverage-report/parse"
	"github.com/xabi93/go-coverage-report/report"
)

func main() {
	in := input.NewLocal(os.Args[1])

	prof, err := in.Load()
	if err != nil {
		log.Fatal(err)
	}

	r, _ := parse.FromLocal(prof)

	fmt.Println(report.Markdown{}.Generate(r))
}
