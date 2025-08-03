package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Реализация запроса к таблице для добавления задачи
func (d *Database) AddTask(task *Task) (int64, error) {

	query := `
		INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)
		`
	res, err := d.DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Реализация запроса к таблице для получения задачи
func (d *Database) GetTask(id string) (*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`

	row := d.DB.QueryRow(query, sql.Named("id", id))

	var task Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &task, nil
}

// Реализация запроса к таблице для изменения задачи
func (d *Database) UpdateTask(task *Task) error {

	query := `UPDATE scheduler  SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := d.DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// Реализация запроса к таблице для удаления задачи
func (d *Database) DeleteTask(id string) error {
	if d.DB == nil {
		return errors.New("database not initialized")
	}

	query := `DELETE FROM scheduler WHERE id = :id`
	res, err := d.DB.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	// Проверяем, что задача была удалена
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if count == 0 {
		return errors.New("error in deliting task")
	}

	return nil
}

// Реализация запроса к таблице для обновления даты
func (d *Database) UpdateDate(id string, next string) error {

	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := d.DB.Exec(query, sql.Named("id", id), sql.Named("date", next))

	if err != nil {
		return fmt.Errorf("update date failed: %w", err)
	}

	// Проверяем, что задача была удалена
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}
	if count == 0 {
		return errors.New("error in deliting task")
	}

	return nil
}
