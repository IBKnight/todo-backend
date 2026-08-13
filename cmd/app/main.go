package main

import (
	"log/slog"
	"os"

	"github.com/IBKnight/todo-backend/internal/app"
)

func main() {
	if err := app.Init(); err != nil {
		slog.Error("failed to run app", "error", err.Error())
		os.Exit(1)
	}
}
