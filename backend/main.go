package main

import (
	"context"
	"fmt"
	"os"
	_ "time/tzdata"

	"github.com/Nesquiko/swimlogs/pkg/server"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
