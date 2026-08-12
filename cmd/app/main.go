package main

import (
	"log/slog"
	"os"
)

func main() {
	if err := Init(); err != nil {
		slog.Error("failed to run app", "error", err.Error())
		os.Exit(1)
	}
}
