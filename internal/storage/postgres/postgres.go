package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	"avito_service/internal/storage"
)

type Storage struct {
	db *sql.DB
}

// / Сделать сегменты уникальными
const (
	CreateSegm = "postgres.CreateSegment"
	DeleteSegm = "postgres.DeleteSegment"
)

var (
	usersQuery = `CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    user_name varchar(100),
    isActual boolean
	)`
	segmQuery = `CREATE TABLE IF NOT EXISTS segments (
		    id serial primary key,
		    slug varchar(100) not null ,
			isActual boolean
		)`

	userSegmQuery = `CREATE TABLE IF NOT EXISTS users_segments(
		    id serial primary key,
		    user_id int references users(id),
		    segment_id int references segments(id),
		    isActive boolean not null,
		    created_at timestamp,
		    deleted_at timestamp,
    		ttl timestamp
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
	query := "INSERT INTO segments (slug, isActual) VALUES ($1,$2)"
	_, err := s.db.Exec(query, slug, true)
	if err != nil {
		return fmt.Errorf("%s: %w", CreateSegm, err)
	}

	return nil
}

func (s *Storage) DeleteSegment(slug string) error {
	query := "UPDATE segments SET isActual = false WHERE slug = $1"

	_, err := s.db.Exec(query, slug)
	if err != nil {
		return fmt.Errorf("%s: %w", DeleteSegm, err)
	}

	query = "UPDATE users_segments SET isActive = false WHERE segment_id IN (SELECT id FROM segments WHERE slug = $1)"

	_, err = s.db.Exec(query, slug)
	if err != nil {
		return fmt.Errorf("%s: %w", DeleteSegm, err)
	}
	return nil
}

func (s *Storage) AddUserToSeg(list []string, userid int, ttl time.Time) error {
	for _, slug := range list {
		queryExists := `
            SELECT COUNT(*) FROM users_segments
            WHERE user_id = $1 AND isActive = true AND segment_id IN (SELECT id FROM segments WHERE slug = $2 AND isActual = true)
        `
		var count int
		err := s.db.QueryRow(queryExists, userid, slug).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			return storage.SegmentAlreadyExistsForUserError{Slug: slug}
		}

		query := "INSERT INTO users_segments (user_id, segment_id, isActive, created_at, ttl) VALUES ($1,(SELECT id FROM segments WHERE slug = $2 AND isActual = true),true, $3, $4)"

		_, err = s.db.Exec(query, userid, slug, time.Now(), ttl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) DeleteSegmentsOfUser(list []string, userid int) error {
	query := "UPDATE users_segments SET isActive = false, deleted_at = $1 WHERE user_id = $2 AND isActive = true AND segment_id IN(SELECT id FROM segments WHERE slug = $3 AND isActual = true)"

	for _, slug := range list {
		_, err := s.db.Exec(query, time.Now(), userid, slug)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) GetActiveSegments(userid int) ([]string, error) {
	query := "SELECT slug FROM segments LEFT JOIN users_segments ON segment_id = segments.id WHERE users_segments.user_id = $1 AND isActive = true"

	var slugs []string
	rows, err := s.db.Query(query, userid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var slug string
		err := rows.Scan(&slug)
		if err != nil {
			return nil, err
		}
		slugs = append(slugs, slug)
	}

	if len(slugs) == 0 {
		return nil, storage.SegmentsNotFound{}
	}
	return slugs, nil
}

func (s *Storage) IfSlugExists(slug string) error {
	query := "SELECT EXISTS (SELECT 1 FROM segments WHERE slug = $1)"

	var exists bool
	_ = s.db.QueryRow(query, slug).Scan(&exists)
	if exists {
		return storage.SegmentAlreadyExistsError{}
	}

	return nil
}
func (s *Storage) IfExists(userId int, slugList []string) error {
	query := "SELECT COUNT(*) FROM users WHERE users.id = $1"
	var count int

	err := s.db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return err
	}

	if count != 1 {
		return storage.UserNotExists{Id: userId}
	}

	query2 := "SELECT EXISTS (SELECT 1 FROM segments WHERE slug = $1 and isActual = true)"
	var exists bool
	for _, slug := range slugList {
		_ = s.db.QueryRow(query2, slug).Scan(&exists)
		if !exists {
			return storage.SegmentNotExists{Slug: slug}
		}
	}

	return nil
}

func (s *Storage) GetReport(userId int, startDate, endDate time.Time) ([]storage.Segment, error) {
	query := `
		SELECT user_id, isActive, created_at, deleted_at
		FROM users_segments
		WHERE user_id = $1
			AND (created_at >= $2 AND created_at <= $3)
			AND (deleted_at >= $2 AND deleted_at <= $3 OR deleted_at IS NULL)
	`
	query2 := "SELECT COUNT(*) FROM users WHERE users.id = $1"
	var count int

	err := s.db.QueryRow(query2, userId).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count != 1 {
		return nil, storage.UserNotExists{Id: userId}
	}
	rows, err := s.db.Query(query, userId, startDate, endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.UserNotExists{}
		}
	}
	defer rows.Close()

	var segments []storage.Segment

	for rows.Next() {
		var segment storage.Segment

		if err := rows.Scan(&segment.UserId, &segment.IsActive, &segment.CreatedAt, &segment.DeletedAt); err != nil {
			return nil, err
		}
		segments = append(segments, segment)
	}

	if len(segments) == 0 {
		return segments, storage.CSVError{}
	}

	return segments, nil
}

func (s *Storage) CheckForTTL() (int, error) {

	current := time.Now().Add(-time.Nanosecond)
	fmt.Println(current)
	query := "UPDATE users_segments SET isActive = false, deleted_at = $1 WHERE ttl <= $2 and isActive = true RETURNING user_id"

	var userId int
	err := s.db.QueryRow(query, current, current).Scan(&userId)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
