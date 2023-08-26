package postgres

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

const (
	CreateSegm = "postgres.CreateSegment"
	DeleteSegm = "postgres.DeleteSegment"
)

var (
	usersQuery = `CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    user_name varchar(100)
	)`
	segmQuery = `CREATE TABLE IF NOT EXISTS segments (
		    id serial primary key,
		    slug varchar(100) not null
		)`

	userSegmQuery = `CREATE TABLE IF NOT EXISTS users_segments(
		    id serial primary key,
		    user_id int references users(id),
		    segment_id int references segments(id),
		    isActive boolean not null,
		    created_at timestamp,
		    deleted_at timestamp
		)`
)

func NewStorage() (*Storage, error) {
	path := os.Getenv("POSTGRES")
	fmt.Println(path)
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	queries := []string{usersQuery, segmQuery, userSegmQuery}

	for _, query := range queries {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, err
		}

		_, err = stmt.Exec()
		if err != nil {
			return nil, err
		}

	}
	return &Storage{db: db}, nil

}

func (s *Storage) CreateUser(name string) error {
	query := "INSERT INTO users (user_name) VALUES ($1)"

	_, err := s.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("%s: %w", CreateSegm, err)
	}

	return nil
}

func (s *Storage) CreateSegment(slug string) error {
	query := "INSERT INTO segments (slug) VALUES ($1)"
	_, err := s.db.Exec(query, slug)
	if err != nil {
		return fmt.Errorf("%s: %w", CreateSegm, err)
	}

	return nil
}

func (s *Storage) DeleteSegment(slug string) error {
	query := "DELETE FROM segments WHERE slug = $1"

	_, err := s.db.Exec(query, slug)
	if err != nil {
		return fmt.Errorf("%s: %w", DeleteSegm, err)
	}
	return nil
}

func (s *Storage) AddUserToSeg(list []string, Userid int) error {
	query := "INSERT INTO users_segments (user_id, segment_id, isActive) VALUES ($1,(SELECT id FROM segments WHERE slug = $2),true)"
	for _, slug := range list {
		_, err := s.db.Exec(query, Userid, slug)
		if err != nil {
			return fmt.Errorf("%s: %w", "AddingToSeg", err)
		}
	}
	return nil
}

//func (s *Storage) DeleteSegmentsOfUser() error {
//
//}
//
//func (s *Storage) GetActiveSegments(id int) []string {
//
//}
