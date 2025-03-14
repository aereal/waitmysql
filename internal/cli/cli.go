package cli

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/aereal/waitmysql/internal/logging"
	"github.com/go-sql-driver/mysql"
	"github.com/shogo82148/go-retry"
)

func New(input io.Reader, outStream, errStream io.Writer) *App {
	return &App{
		input:     input,
		outStream: outStream,
		errStream: errStream,
	}
}

type App struct {
	input                io.Reader
	outStream, errStream io.Writer
}

const (
	StatusOK    = 0
	StatusError = 1
)

func (a *App) Run(ctx context.Context, argv []string) int {
	if err := a.run(ctx, argv); err != nil {
		logging.Error(ctx, err)
		return StatusError
	}
	return StatusOK
}

func (a *App) run(ctx context.Context, argv []string) error {
	flagset := flag.NewFlagSet(argv[0], flag.ContinueOnError)
	var (
		dsn         string
		initialWait time.Duration
		maxDelay    time.Duration
		maxAttempts int
	)
	flagset.StringVar(&dsn, "dsn", "", "data source name string")
	flagset.DurationVar(&initialWait, "init-wait", 0, "initial delay")
	flagset.DurationVar(&maxDelay, "max-delay", -1, "max delay (default is infinity)")
	flagset.IntVar(&maxAttempts, "max-attempts", -1, "max attempts count")
	if err := flagset.Parse(argv[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if dsn == "" {
		return ErrMissingDSN
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return fmt.Errorf("mysql.ParseDSN: %w", err)
	}
	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		return fmt.Errorf("mysql.NewConnector: %w", err)
	}
	db := sql.OpenDB(connector)
	checkFn := func() error {
		return db.PingContext(ctx)
	}
	policy := &retry.Policy{
		MinDelay: initialWait,
		MaxDelay: maxDelay,
		MaxCount: maxAttempts,
	}
	if err := policy.Do(ctx, checkFn); err != nil {
		return err
	}
	return nil
}
