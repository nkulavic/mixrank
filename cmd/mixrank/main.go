package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nkulavic/mixrank/internal/cli"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cmd := cli.New()
	if e := cmd.ExecuteContext(ctx); e != nil {
		b, _ := json.Marshal(map[string]any{"error": e.Error(), "exit_code": cli.ExitCode(e)})
		fmt.Fprintln(os.Stderr, string(b))
		os.Exit(cli.ExitCode(e))
	}
}
