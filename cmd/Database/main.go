package main

import (
	"context"
	"flag"
	"log/slog"
	"log_parser/database"
	"log_parser/segmenter"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	logPath := flag.String("path", "/home/navaneeth/project_log_parser/logs", "Path to the log directory")
	flag.Parse()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	logStore, err := segmenter.GenerateSegments(*logPath)
	if err != nil {
		slog.Error("Failed to parse logs", "error", err)

	}

	err = database.InsertLogs(ctx, conn, logStore)
	if err != nil {
		slog.Error("Failed to insert logs", "error", err)
	} else {
		slog.Info("All logs inserted successfully")
	}
}
