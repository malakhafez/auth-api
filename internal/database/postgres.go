
package database

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	// Read settings from the .env file.
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("could not load .env file: %w", err)
	}

	// Read database settings.
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Check that the required settings exist.
	if host == "" || port == "" || user == "" ||
		password == "" || dbName == "" {
		return nil, fmt.Errorf("database settings are incomplete")
	}

	// Build the connection URL safely.
	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   dbName,
	}
	query := connectionURL.Query()
	query.Set("sslmode", "disable")
	connectionURL.RawQuery = query.Encode()

	// Open the database connection.
	db, err := sql.Open("pgx", connectionURL.String())
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	// Check whether the database is reachable.
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return db, nil
}
