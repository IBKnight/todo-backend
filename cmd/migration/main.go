package main

// import (
// 	"flag"
// 	"log/slog"
// 	"os"
// )

// func main() {
// 	direction := flag.String("direction", "up", "up or down")
// 	flag.Parse()

// 	db, err := connectDB()
// 	if err != nil {
// 		slog.Error("db connect", "error", err)
// 		os.Exit(1)
// 	}

// 	switch *direction {
// 	case "up":
// 		err = migrations.Up(db)
// 	case "down":
// 		err = migrations.Down(db)
// 	}

// 	if err != nil {
// 		slog.Error("migration failed", "error", err)
// 		os.Exit(1)
// 	}
// }
