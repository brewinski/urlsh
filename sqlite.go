package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./urlshortener.db?_debug=1")
	if err != nil {
		return nil, fmt.Errorf("error opening new sqlite db connection, %w", err)
	}

	// SQL statement to create the todos table if it doesn't exist
	sqlStmt := `
	 CREATE TABLE IF NOT EXISTS short_urls(
	  ID TEXT NOT NULL PRIMARY KEY,
	  URL TEXT
	 );`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return db, fmt.Errorf("error creating table: %w: %s\n", err, sqlStmt) // Log an error if table creation fails
	}

	return db, nil
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db,
	}
}

type ShortURL struct {
	ID  string
	URL string
}

func (store *Store) Get(id string) (ShortURL, error) {
	// Add debug logging
	fmt.Printf("Searching for ID: %q\n", id)

	// Use QueryRow instead of Query since we're looking for a single result
	query := "SELECT ID, URL FROM short_urls WHERE ID = ?"
	row := store.db.QueryRow(query, id)

	var url ShortURL
	err := row.Scan(&url.ID, &url.URL)
	if err == sql.ErrNoRows {
		// No result found
		return ShortURL{}, nil
	}
	if err != nil {
		return ShortURL{}, fmt.Errorf("get by short url failed: %w", err)
	}

	// Debug: print what we found
	fmt.Printf("Found URL: %+v\n", url)

	return url, nil
}
func (store *Store) Set(id, url string) error {
	query := `INSERT INTO short_urls (ID, URL) VALUES (?, ?)`

	result, err := store.db.Exec(query, id, url)
	if err != nil {
		return fmt.Errorf("CreateShortURL: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CreateShortURL rows affected: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("CreateShortURL: expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (store *Store) ListShortURL() ([]ShortURL, error) {
	rows, err := store.db.Query("SELECT * FROM short_urls")
	if err != nil {
		return []ShortURL{}, fmt.Errorf("ListShortURL: %w", err)
	}

	defer rows.Close()

	urls := []ShortURL{}
	for rows.Next() {
		var url ShortURL
		if err := rows.Scan(&url.ID, &url.URL); err != nil {
			return nil, fmt.Errorf("ListShortURL scan: %w", err)
		}

		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ListShortURL iteration: %w", err)
	}

	return urls, nil

}

func (store *Store) CreateShortURL(shortUrl ShortURL) error {
	query := `INSERT INTO short_urls (ID, URL) VALUES (?, ?)`

	result, err := store.db.Exec(query, shortUrl.ID, shortUrl.URL)
	if err != nil {
		return fmt.Errorf("CreateShortURL: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CreateShortURL rows affected: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("CreateShortURL: expected to affect 1 row, affected %d", rows)
	}

	return nil
}
