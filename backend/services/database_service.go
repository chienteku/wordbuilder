package services

import (
	"database/sql"
	"time"

	"wordbuilder/models"

	_ "github.com/mattn/go-sqlite3"
)

// DatabaseService handles database operations
type DatabaseService struct {
	DB *sql.DB
}

// NewDatabaseService creates a new database service
func NewDatabaseService(dbPath string) (*DatabaseService, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	service := &DatabaseService{DB: db}
	err = service.InitTables()
	if err != nil {
		return nil, err
	}

	return service, nil
}

// InitTables creates necessary tables if they don't exist
func (s *DatabaseService) InitTables() error {
	wordListTable := `
	CREATE TABLE IF NOT EXISTS word_lists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		source TEXT,
		file_path TEXT NOT NULL,
		word_count INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`

	_, err := s.DB.Exec(wordListTable)
	return err
}

// InsertWordList adds a new word list to the database
func (s *DatabaseService) InsertWordList(list *models.WordList) (int, error) {
	now := time.Now()
	list.CreatedAt = now
	list.UpdatedAt = now

	result, err := s.DB.Exec(
		"INSERT INTO word_lists (name, description, source, file_path, word_count, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		list.Name, list.Description, list.Source, list.FilePath, list.WordCount, list.CreatedAt, list.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	list.ID = int(id)
	return list.ID, nil
}

// GetWordList retrieves a word list by ID
func (s *DatabaseService) GetWordList(id int) (*models.WordList, error) {
	var list models.WordList
	err := s.DB.QueryRow(
		"SELECT id, name, description, source, file_path, word_count, created_at, updated_at FROM word_lists WHERE id = ?",
		id,
	).Scan(&list.ID, &list.Name, &list.Description, &list.Source, &list.FilePath, &list.WordCount, &list.CreatedAt, &list.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &list, nil
}

// GetAllWordLists retrieves all word lists
func (s *DatabaseService) GetAllWordLists() ([]*models.WordList, error) {
	rows, err := s.DB.Query("SELECT id, name, description, source, file_path, word_count, created_at, updated_at FROM word_lists ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []*models.WordList
	for rows.Next() {
		var list models.WordList
		if err := rows.Scan(&list.ID, &list.Name, &list.Description, &list.Source, &list.FilePath, &list.WordCount, &list.CreatedAt, &list.UpdatedAt); err != nil {
			return nil, err
		}
		lists = append(lists, &list)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

// UpdateWordList updates an existing word list
func (s *DatabaseService) UpdateWordList(list *models.WordList) error {
	list.UpdatedAt = time.Now()

	_, err := s.DB.Exec(
		"UPDATE word_lists SET name = ?, description = ?, source = ?, file_path = ?, word_count = ?, updated_at = ? WHERE id = ?",
		list.Name, list.Description, list.Source, list.FilePath, list.WordCount, list.UpdatedAt, list.ID,
	)

	return err
}

// DeleteWordList removes a word list by ID
func (s *DatabaseService) DeleteWordList(id int) error {
	_, err := s.DB.Exec("DELETE FROM word_lists WHERE id = ?", id)
	return err
}

// Close closes the database connection
func (s *DatabaseService) Close() error {
	return s.DB.Close()
}
