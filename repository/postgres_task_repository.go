package repository

import (
	"context"
	"database/sql"
	"errors"

	"example.com/go-backend/model"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	query := `
		SELECT id, title, is_done, COALESCE(description, ''), priority
		FROM tasks
		ORDER BY id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.IsDone,
			&task.Description,
			&task.Priority,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task model.Task) (model.Task, error) {
	query := `
		INSERT INTO tasks (title, is_done, description, priority)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, is_done, COALESCE(description, ''), priority;
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		task.Title,
		task.IsDone,
		task.Description,
		task.Priority,
	)

	var createdTask model.Task

	if err := row.Scan(
		&createdTask.ID,
		&createdTask.Title,
		&createdTask.IsDone,
		&createdTask.Description,
		&createdTask.Priority,
	); err != nil {
		return model.Task{}, err
	}

	return createdTask, nil
}

func (r *PostgresTaskRepository) Update(ctx context.Context, id int, task model.Task) (model.Task, error) {
	query := `
		UPDATE tasks
		SET title = $1,
			is_done = $2,
			description = $3,
			priority = $4
		WHERE id = $5
		RETURNING id, title, is_done, COALESCE(description, ''), priority;
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		task.Title,
		task.IsDone,
		task.Description,
		task.Priority,
		id,
	)

	var updatedTask model.Task

	err := row.Scan(
		&updatedTask.ID,
		&updatedTask.Title,
		&updatedTask.IsDone,
		&updatedTask.Description,
		&updatedTask.Priority,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return updatedTask, nil
}

func (r *PostgresTaskRepository) Patch(ctx context.Context, id int, patch model.TaskPatch) (model.Task, error) {
	query := `
		UPDATE tasks
		SET title = COALESCE($1, title),
			is_done = COALESCE($2, is_done),
			description = COALESCE($3, description),
			priority = COALESCE($4, priority)
		WHERE id = $5
		RETURNING id, title, is_done, COALESCE(description, ''), priority;
	`

	row := r.db.QueryRowContext(
		ctx,
		query,
		patch.Title,
		patch.IsDone,
		patch.Description,
		patch.Priority,
		id,
	)

	var patchedTask model.Task

	err := row.Scan(
		&patchedTask.ID,
		&patchedTask.Title,
		&patchedTask.IsDone,
		&patchedTask.Description,
		&patchedTask.Priority,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return patchedTask, nil
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

var _ TaskRepository = (*PostgresTaskRepository)(nil)
