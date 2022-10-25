package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/xabi93/go-coverage-report/pkg/cover"
	"github.com/xabi93/go-coverage-report/pkg/format"
	"github.com/xabi93/go-coverage-report/pkg/notify/github"
)

var (
	errMissingOwnerOrRepo     = errors.New("missing owner or repository")
	errOnlyGHActionRunSupport = errors.New("Only supports running in github actions")
)

const name = "go-coverage-report"

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	flags, opts := setupFlags(name)

	switch err := flags.Parse(args); {
	case err == pflag.ErrHelp:
		return nil
	case err != nil:
		usage(os.Stderr, name, flags)
		return err
	}

	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return errOnlyGHActionRunSupport
	}

	notifier, err := loadGHActions(ctx, opts)
	if err != nil {
		return err
	}

	formatter, err := format.NewMarkdown(opts.Template)
	if err != nil {
		return err
	}

	report, err := cover.NewFileParser(opts.CoverageFile).Parse(ctx)
	if err != nil {
		return err
	}

	out, err := formatter.Format(report)
	if err != nil {
		return err
	}

	return notifier.Notify(ctx, out)
}

func loadGHActions(ctx context.Context, opts *options) (*github.Notifier, error) {
	ownerRepo := strings.Split(opts.GHRepository, "/")

	if len(ownerRepo) < 2 {
		return nil, errMissingOwnerOrRepo
	}

	return github.NewNotifier(
		github.NewClient(ctx, opts.GHToken).Checks,
		ownerRepo[0],
		ownerRepo[1],
		opts.GHSha,
		opts.ReportName,
	), nil
}

type options struct {
	ReportName   string
	CoverageFile string

	GHToken      string
	GHRepository string
	GHSha        string

	Template string
}

func setupFlags(name string) (*pflag.FlagSet, *options) {
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)
	flags.SetInterspersed(false)
	flags.Usage = func() {
		usage(os.Stdout, name, flags)
	}

	opt := &options{}

	flags.StringVar(&opt.CoverageFile, "coverage-file", "", "Path where the coverage file is located.")
	flags.StringVar(&opt.ReportName, "report-name", "", "Title of the coverage report")
	flags.StringVar(&opt.GHToken, "github-token", os.Getenv("GITHUB_TOKEN"), "Github authentication token. (env: GITHUB_TOKEN)")
	flags.StringVar(&opt.GHRepository, "github-repository", os.Getenv("GITHUB_REPOSITORY"), "Repository name with owner. ex: octocat/Hello-World ex: octocat")
	flags.StringVar(&opt.GHSha, "github-sha", os.Getenv("GITHUB_SHA"), "The commit SHA that triggered the workflow. ex: ffac537e6cbbf934b08745a378932722df287a53")
	flags.StringVar(&opt.Template, "template", "", "Custom template for output")

	return flags, opt
}

func usage(out io.Writer, name string, flags *pflag.FlagSet) {
	fmt.Fprintf(out, `Usage:
    %[1]s [flags]
    %[1]s
Flags:
`, name)
	flags.SetOutput(out)
	flags.PrintDefaults()
}
