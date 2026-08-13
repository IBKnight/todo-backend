package main

import (
	"github.com/IBKnight/todo-backend/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := app.Init(); err != nil {
		logrus.Fatalf("failed to run app: %v", err.Error())
	}
}
