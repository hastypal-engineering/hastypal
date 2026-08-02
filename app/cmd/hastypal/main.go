package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/adriein/hastypal/internal"
	"github.com/adriein/hastypal/internal/server"
)

func main() {
	app := internal.NewApp()

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.Shutdown(ctx); err != nil {
			fmt.Printf("Failed to perform graceful shutdown %v", err)
			os.Exit(1)
		}
	}()

	if len(os.Args) < 2 {
		server.New(app)

		return
	}

	switch os.Args[1] {
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

}
