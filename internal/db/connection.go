package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/JamesAndresCM/eco_orders/pkg/logger"
)

var DB *sql.DB

func InitDB() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST") // ej: "localhost"
	dbPort := os.Getenv("DB_PORT") // ej: "5432"

	var dbName string
	if env == "production" {
		dbName = "eco"
	} else {
		dbName = "eco_development"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Error opening DB connection:", err)
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		logger.Error("Error pinging DB:", err)
		panic(err)
	}

	logger.Success(fmt.Sprintf("Connected to database %s in %s mode", dbName, env))
}

