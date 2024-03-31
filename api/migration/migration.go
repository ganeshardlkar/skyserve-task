package migration

import (
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

func ExecuteMigration(db *sqlx.DB, path string) error {
	migrationSQL, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error reding migration file")
		return err
	}

	statements := strings.Split(string(migrationSQL), ";")

	for _, statement := range statements {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func ExecuteMigrationUp(db *sqlx.DB, path string) error {
	migrationSQL, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error while up migration")
		return nil
	}

	statements := strings.Split(string(migrationSQL), ";")
	for _, statement := range statements {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		_, err := db.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExecuteMigrationDown(db *sqlx.DB, migrationFilePath string) error {
	migrationSQL, err := os.ReadFile(migrationFilePath)
	if err != nil {
		return err
	}

	statements := strings.Split(string(migrationSQL), ";")
	for i := len(statements) - 1; i >= 0; i-- {
		if strings.TrimSpace(statements[i]) == "" {
			continue
		}
		_, err := db.Exec(statements[i])
		if err != nil {
			return err
		}
	}

	return nil
}
