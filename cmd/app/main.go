package main

import (
	"log"

	todobackend "github.com/IBKnight/todo-backend"
	"github.com/IBKnight/todo-backend/pkg/handler"
)

func main() {
	handler := handler.Handler{}

	srv := new(todobackend.Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
