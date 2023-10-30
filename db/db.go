package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/ralphvw/sprint-planner-v2/helpers"
  "github.com/golang-migrate/migrate/v4"
)

const (
  migrationsDir = "file://./db/migrations"
)

func InitDb() *sql.DB {
  var err error

  db, err := sql.Open("postgres", os.Getenv("SPRINT_DB_URL"))

  if err != nil {
    log.Fatal(err)
  }

  helpers.LogAction("DB Connected successfully")

  if err := applyMigrations(); err != nil {
    log.Fatal(err)
  }

  return db
}

func applyMigrations() error {
  m, err := migrate.New(migrationsDir, os.Getenv("SPRINT_DB_URL"))
  if err != nil {
  return err
  }

  if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    return err
  }

  return nil;
}
