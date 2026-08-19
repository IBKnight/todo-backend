package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBKnight/todo-backend/internal/handler"
	"github.com/IBKnight/todo-backend/internal/repository"
	"github.com/IBKnight/todo-backend/internal/service/auth"
	todoitem "github.com/IBKnight/todo-backend/internal/service/todo_item"
	todolist "github.com/IBKnight/todo-backend/internal/service/todo_list"
	"github.com/IBKnight/todo-backend/migrations"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error occured while env init: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		return fmt.Errorf("error occured while configs init: %s", err.Error())
	}

	port := viper.GetString("port")

	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbUsername := viper.GetString("db.username")
	dbName := viper.GetString("db.dbname")
	dbSSLMode := viper.GetString("db.sslmode")

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return errors.New("DB_PASSWORD env var is required")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		return errors.New("SECRET env var is required")
	}

	tokenTTL := time.Hour

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
		return fmt.Errorf("error occured while db init: %w", err)
	}

	if err := migrations.Up(db); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	authRepo := repository.NewAuthRepo(db)
	todoListRepo := repository.NewTodoListRepo(db)
	todoItemRepo := repository.NewTodoItemRepo(db)

	authService := auth.NewService(authRepo, []byte(secret), tokenTTL)
	todolistService := todolist.NewService(todoListRepo)
	todoItemService := todoitem.NewService(todoItemRepo)

	h := handler.NewHandler(
		authService,
		todolistService,
		todoItemService,
	)

	srv := NewServer(port, h.InitRoutes())

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("error occured while running http server: %s", err)
		}
	}()

	logrus.Info("todo app started")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logrus.Info("todo app shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err)
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err)
	}

	logrus.Info("todo app stopped")

	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
