package main

import (
	"fmt"
	"log/slog"
	"log_parser/pkg/models"
	"log_parser/pkg/parser"
	"os"
)

// const dbUrl = "postgresql:///log_analyser_db?host=/var/run/postgresql/"
const dbUrl = "postgresql:///log_data?host=/var/run/postgresql/"

func handleCommand(args []string) error {

	db, err := models.CreateDB(dbUrl)
	if err != nil {
		return err
	}
	switch args[0] {

	case "init":
		err := models.InitDb(db)
		if err != nil {
			return err
		}

	case "add":
		fname := args[1] // TBD : Handle case where no filename specified
		if fname == "" {
			slog.Error("Specify directory!")
		}
		entries, err := parser.ParseLogFiles(fname)
		if err != nil {
			return err
		}

		for _, e := range entries {

			var level models.LogLevel
			if err := db.First(&level, "level = ?", string(e.Level)).Error; err != nil {
				return fmt.Errorf("unknown level %s: %w", e.Level, err)
			}

			//  Look up ComponentID
			var component models.LogComponent
			if err := db.First(&component, "component = ?", e.Component).Error; err != nil {
				return fmt.Errorf("unknown component %s: %w", e.Component, err)
			}

			// models.AddEntry(db, e)

			//  Look up HostID
			var host models.LogHost
			if err := db.First(&host, "host = ?", e.Host).Error; err != nil {
				return fmt.Errorf("unknown host %s: %w", e.Host, err)
			}

			//  Create Entry WITH foreign keys
			dbEntry := models.Entry{
				TimeStamp:   e.Time,
				LevelID:     uint(level.ID),
				ComponentID: component.ID,
				HostID:      host.ID,
				RequestID:   e.Request_id,
				Message:     e.Message,
			}

			if err := models.AddDb(db, dbEntry); err != nil {
				return err
			}
		}
		return nil

	case "query":
		queries := args[1:]
		entries, err := models.QueryDB(db, queries)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fmt.Println(entry)
		}
		slog.Info("Filtering successful!", "no. of entries:", len(entries))
		return nil

	default:
		slog.Warn("Unknown command!")
		return fmt.Errorf("uknown command %v", args[0])
	}
	return nil
}

func main() {
	err := handleCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in invocation %v", err)
		os.Exit(-1)
	}
}
