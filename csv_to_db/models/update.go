package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the database connection
func InitDB() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get PostgreSQL connection details
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Build the connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Calcutta",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Open a connection to the database
	var err2 error
	DB, err2 = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err2 != nil {
		log.Fatalf("Error opening connection to database: %v", err2)
	}

	// Migrate the schema (create table based on struct)
	err2 = DB.AutoMigrate(&Update{})
	if err2 != nil {
		log.Fatalf("Error migrating schema: %v", err2)
	}
}

type Update struct {
	ID                    uint   `gorm:"primaryKey"`
	FxiletID              string `gorm:"type:varchar(255);not null"`
	Name                  string `gorm:"type:text;not null"`
	Criticality           string `gorm:"type:varchar(50);not null"`
	RelevantComputerCount int    `gorm:"type:int;not null"`
}
