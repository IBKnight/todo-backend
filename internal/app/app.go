package app

import (
	"fmt"
	"os"

	todobackend "github.com/IBKnight/todo-backend"
	"github.com/IBKnight/todo-backend/internal/handler"
	"github.com/IBKnight/todo-backend/internal/repository"
	"github.com/IBKnight/todo-backend/internal/service/auth"
	todoitem "github.com/IBKnight/todo-backend/internal/service/todo_item"
	todolist "github.com/IBKnight/todo-backend/internal/service/todo_list"
	"github.com/IBKnight/todo-backend/migrations"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() error {
	if err := initConfig(); err != nil {
		return fmt.Errorf("error occured while configs init: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error occured while env init: %s", err.Error())
	}

	port := viper.GetString("port")

	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbUsername := viper.GetString("db.username")
	dbName := viper.GetString("db.dbname")
	dbSSLMode := viper.GetString("db.sslmode")
	dbPassword := os.Getenv("DB_PASSWORD")

	db, err := repository.NewPostgresDB(
		&repository.Config{
			Host:     dbHost,
			Port:     dbPort,
			Username: dbUsername,
			Password: dbPassword,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		})

	if err != nil {
		return fmt.Errorf("error occured while db init: %s", err.Error())

	}

	if err := migrations.Up(db); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	authRepo := repository.NewAuthRepo(db)
	todoListRepo := repository.NewTodoListRepo(db)
	todoItemRepo := repository.NewTodoItemRepo(db)

	authService := auth.NewService(authRepo)
	todolistService := todolist.NewService(todoListRepo)
	todoItemService := todoitem.NewService(todoItemRepo)

	handler := handler.NewHandler(
		authService,
		todolistService,
		todoItemService,
	)

	srv := new(todobackend.Server)
	if err := srv.Run(port, handler.InitRoutes()); err != nil {
		return fmt.Errorf("error occured while running http server: %s", err.Error())
	}

	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
