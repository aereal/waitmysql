package main

import (
	"context"
	"os"

	"github.com/aereal/waitmysql/internal/cli"
)

func main() {
	app := cli.New(os.Stdin, os.Stdout, os.Stderr)
	os.Exit(app.Run(context.Background(), os.Args))
}
