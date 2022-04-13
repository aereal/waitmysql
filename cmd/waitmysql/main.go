package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shogo82148/go-retry"
)

func main() {
	os.Exit(int(run(os.Args)))
}

type statusCode int

const (
	statusOK statusCode = iota
	statusError
)

func run(argv []string) statusCode {
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
	switch err := flagset.Parse(argv[1:]); err {
	case nil:
	case flag.ErrHelp:
		return statusOK
	default:
		fmt.Printf("%+v\n", err)
		return statusError
	}
	if dsn == "" {
		fmt.Println("-dsn must be given")
		return statusError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	policy := &retry.Policy{
		MinDelay: initialWait,
		MaxDelay: maxDelay,
		MaxCount: maxAttempts,
	}
	if err := policy.Do(ctx, func() error { return checkConnection(ctx, dsn) }); err != nil {
		fmt.Printf("%+v\n", err)
		return statusError
	}
	return statusOK
}

func checkConnection(ctx context.Context, dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return retry.MarkPermanent(err)
	}
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
