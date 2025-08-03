package db

import (
	"database/sql"
	"fmt"
)

func (d *Database) Tasks(limit int) ([]*Task, error) {
	query := `
	   SELECT id, date, title, comment, repeat
	   FROM scheduler
	   ORDER BY date ASC
	   LIMIT :limit
   `
	rows, err := d.DB.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}
