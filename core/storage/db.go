package storage

import (
	"database/sql"
	"log"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type SearchResult struct {
	id         int
	content    string
	created_at time.Time
}

func InitDb() (*Storage, error) {
	// TODO: Get db path based on the OS
	db, err := sql.Open("sqlite3", path.Join(".", "tmp", "clipboard_history.db"))

	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) search_for_exact_match(content string) ([]SearchResult, error) {
	rows, err := s.db.Query("SELECT * FROM clipboard_history WHERE content = ?", content)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := []SearchResult{}

	for rows.Next() {
		var id int
		var content string
		var created_at time.Time

		if err := rows.Scan(&id, &content, &created_at); err != nil {
			return nil, err
		}

		results = append(results, SearchResult{id, content, created_at})

	}

	return results, rows.Err()
}

func (s *Storage) Add_Clipboard_Entry(content string) error {
	existing_rows, search_err := s.search_for_exact_match(content)

	if search_err != nil {
		log.Fatal("Failed to search for exact match: ", search_err)
		return search_err
	}
	// if there are no existing entries - create a new one
	if len(existing_rows) == 0 {
		displayContent := content
		if len(content) > 15 {
			displayContent = content[:15] + "..."
		}
		log.Println("Adding new entry to db for content: ", displayContent)
		_, err := s.db.Exec("INSERT INTO clipboard_history (content) VALUES (?)", content)
		return err
	} else {
		// if there are existing one - update the first one and delete the rest
		first_row := existing_rows[0]
		rows_to_delete := existing_rows[1:]

		displayContent := first_row.content
		if len(first_row.content) > 15 {
			displayContent = first_row.content[:15] + "..."
		}
		log.Println("Updating entry in db for content: ", displayContent)
		_, err := s.db.Exec("UPDATE clipboard_history SET created_at = ? WHERE id = ?", time.Now(), first_row.id)
		if err != nil {
			log.Fatal("Failed to update entry in db: ", err)
			return err
		}

		for _, row := range rows_to_delete {
			log.Println("Deleting duplicated entry in db with id: ", row.id)
			_, err := s.db.Exec("DELETE FROM clipboard_history WHERE id = ?", row.id)

			if err != nil {
				log.Fatal("Failed to delete duplicated entry in db: ", err)
				return err
			}
		}
	}

	return nil
}
