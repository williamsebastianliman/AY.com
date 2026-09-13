package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
  dsn := os.Getenv("DB_URL")
  if dsn == "" {
    log.Fatal("DB_URL must be set")
  }
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    log.Fatalf("failed to connect to DB: %v", err)
  }
  return db
}