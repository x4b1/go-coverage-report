package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"

	"github.com/xabi93/go-coverage-report/internal/log"
	"github.com/xabi93/go-coverage-report/pkg/cover"
	"github.com/xabi93/go-coverage-report/pkg/format"
	"github.com/xabi93/go-coverage-report/pkg/notify"
	"github.com/xabi93/go-coverage-report/pkg/notify/github"
	"github.com/xabi93/go-coverage-report/pkg/notify/stdout"
)

var (
	errMissingOwnerOrRepo     = errors.New("missing owner or repository")
	errOnlyGHActionRunSupport = errors.New("Only supports running in github actions")
)

const name = "go-coverage-report"

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		log.Error(err.Error())
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

	if opts.debug {
		log.SetLevel(log.DebugLevel)
	}

	var notifier notify.Notifier
	var err error

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		notifier, err = loadGHActions(ctx, opts)
		if err != nil {
			return err
		}
	} else {
		notifier = stdout.NewNotifier()
	}

	formatter, err := format.NewMarkdown(opts.template)
	if err != nil {
		return err
	}

	report, err := cover.NewFileParser(opts.coverageFile).Parse(ctx)
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
	ownerRepo := strings.Split(opts.ghRepository, "/")

	if len(ownerRepo) < 2 {
		return nil, errMissingOwnerOrRepo
	}

	return github.NewNotifier(
		github.NewClient(ctx, opts.ghToken).Checks,
		ownerRepo[0],
		ownerRepo[1],
		opts.ghSha,
		opts.reportName,
	), nil
}

type options struct {
	debug        bool
	reportName   string
	coverageFile string

	ghToken      string
	ghRepository string
	ghSha        string

	template string
}

func setupFlags(name string) (*pflag.FlagSet, *options) {
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)
	flags.SetInterspersed(false)
	flags.Usage = func() {
		usage(os.Stdout, name, flags)
	}

	opts := &options{}

	flags.StringVar(&opts.coverageFile, "coverage-file", "", "Path where the coverage file is located.")
	flags.StringVar(&opts.reportName, "report-name", "", "Title of the coverage report")
	flags.StringVar(&opts.ghToken, "github-token", os.Getenv("GITHUB_TOKEN"), "Github authentication token. (env: GITHUB_TOKEN)")
	flags.StringVar(&opts.ghRepository, "github-repository", os.Getenv("GITHUB_REPOSITORY"), "Repository name with owner. ex: octocat/Hello-World ex: octocat")
	flags.StringVar(&opts.ghSha, "github-sha", os.Getenv("GITHUB_SHA"), "The commit SHA that triggered the workflow. ex: ffac537e6cbbf934b08745a378932722df287a53")
	flags.StringVar(&opts.template, "template", "", "Custom template for output")
	flags.BoolVar(&opts.debug, "debug", false, "enabled debug logging")

	return flags, opts
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
