package database

import (
	"CrispyBot/variables"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Connect establishes a connection to the PostgreSQL database
func Connect() *sql.DB {
	// Open the PostgreSQL database
	db, err := sql.Open("postgres", variables.DB)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %w", err))
	}

	// Initialize the database schema
	if err = initializeSchema(db); err != nil {
		panic(fmt.Errorf("failed to initialize database schema: %w", err))
	}

	return db
}

// initializeSchema creates all necessary tables if they don't exist
func initializeSchema(db *sql.DB) error {
	// Create the users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			discord_id TEXT PRIMARY KEY
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create the characters table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS characters (
			id TEXT PRIMARY KEY,
			owner TEXT NOT NULL,
			FOREIGN KEY(owner) REFERENCES users(discord_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create characters table: %w", err)
	}

	// Create the stats table with a SERIAL primary key for PostgreSQL
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS stats (
			id SERIAL PRIMARY KEY,
			character_id TEXT NOT NULL,
			stat_type INTEGER NOT NULL,
			stat_name TEXT NOT NULL,
			rarity TEXT NOT NULL,
			value INTEGER NOT NULL,
			FOREIGN KEY(character_id) REFERENCES characters(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create stats table: %w", err)
	}

	// Create the traits table with a SERIAL primary key
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS traits (
			id SERIAL PRIMARY KEY,
			character_id TEXT NOT NULL,
			trait_type INTEGER NOT NULL,
			trait_name TEXT NOT NULL,
			rarity TEXT NOT NULL,
			trait_category TEXT NOT NULL,
			FOREIGN KEY(character_id) REFERENCES characters(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create traits table: %w", err)
	}

	// Create the trait_stat_values table with a SERIAL primary key
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS trait_stat_values (
			id SERIAL PRIMARY KEY,
			trait_id INTEGER NOT NULL,
			stat_name TEXT NOT NULL,
			value INTEGER NOT NULL,
			FOREIGN KEY(trait_id) REFERENCES traits(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create trait_stat_values table: %w", err)
	}

	return nil
}
