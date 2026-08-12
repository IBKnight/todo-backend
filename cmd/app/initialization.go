package main

import (
	"fmt"

	todobackend "github.com/IBKnight/todo-backend"
	"github.com/IBKnight/todo-backend/internal/handler"
	"github.com/IBKnight/todo-backend/internal/repository"
	"github.com/IBKnight/todo-backend/internal/service/auth"
	todoitem "github.com/IBKnight/todo-backend/internal/service/todo_item"
	todolist "github.com/IBKnight/todo-backend/internal/service/todo_list"
)

func Init() error {
	authRepo := repository.NewAuthRepo()
	todoListRepo := repository.NewTodoListRepo()
	todoItemRepo := repository.NewTodoItemRepo()

	authService := auth.NewService(authRepo)
	todolistService := todolist.NewService(todoListRepo)
	todoItemService := todoitem.NewService(todoItemRepo)

	handler := handler.NewHandler(
		authService,
		todolistService,
		todoItemService,
	)

	srv := new(todobackend.Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		return fmt.Errorf("error occured while running http server: %s", err.Error())
	}

	return nil
}
