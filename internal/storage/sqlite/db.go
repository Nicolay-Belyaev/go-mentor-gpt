package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"go-mentor-gpt/domain"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

type StorageSQLite struct {
	db *sql.DB
}

func NewStorageSQLite(path string) (*StorageSQLite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite failed: %w", err)
	}

	s := &StorageSQLite{db: db}

	if err := s.initScheme(); err != nil {
		db.Close()
		return nil, fmt.Errorf("init scheme failed: %w", err)
	}

	return s, nil
}

func (s *StorageSQLite) initScheme() error {
	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageSQLite) Close() error {
	return s.db.Close()
}

// Topic Repo
func (s *StorageSQLite) GetAllTopics(ctx context.Context) ([]domain.Topic, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, key, title, comment, status, last_reviewed_at FROM topics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic
	for rows.Next() {
		var t domain.Topic
		if err := rows.Scan(
			&t.ID,
			&t.Key,
			&t.Title,
			&t.Comment,
			&t.Status,
			&t.LastReviewed); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return topics, nil
}

func (s *StorageSQLite) CreateTopic(ctx context.Context, topic domain.Topic) (int64, error) {
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO topics (key, title, comment, status, last_reviewed_at)
			   VALUES (?, ?, ?, ?, ?)`,
		topic.Key,
		topic.Title,
		topic.Comment,
		topic.Status,
		topic.LastReviewed)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *StorageSQLite) UpdateTopic(ctx context.Context, topic domain.Topic) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE topics SET comment = ?, status = ?, last_reviewed_at = ? WHERE id = ?`,
		topic.Comment,
		topic.Status,
		topic.LastReviewed,
		topic.ID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Task Repo
func (s *StorageSQLite) GetAllTasks(ctx context.Context) ([]domain.Task, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, title, comment, status, last_reviewed_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Comment,
			&t.Status,
			&t.LastReviewed); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *StorageSQLite) CreateTask(ctx context.Context, task domain.Task) (int64, error) {
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO tasks (title, comment, status, last_reviewed_at)
			   VALUES (?, ?, ?, ?)`,
		task.Title,
		task.Comment,
		task.Status,
		task.LastReviewed)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *StorageSQLite) UpdateTask(ctx context.Context, task domain.Task) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE tasks SET comment = ?, status = ?, last_reviewed_at = ? WHERE id = ?`,
		task.Comment,
		task.Status,
		task.LastReviewed,
		task.ID)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Goal Repo

func (s *StorageSQLite) GetGoal(ctx context.Context) (*domain.Goal, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, text, comment FROM goal WHERE id = 1")
	var goal domain.Goal
	if err := row.Scan(&goal.ID, &goal.Text, &goal.Comment); err != nil {
		return nil, err
	}
	return &goal, nil
}

func (s *StorageSQLite) SetGoal(ctx context.Context, goal *domain.Goal) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE goal
		 SET text = ?, comment = ?
		 WHERE id = 1`,
		goal.Text,
		goal.Comment,
	)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}
