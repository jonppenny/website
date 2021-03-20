package mysql

import (
	"database/sql"
	"errors"
	"jonppenny.co.uk/webapp/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content, status string) (int, error) {
	stmt := `INSERT INTO posts (title, content, status, created, updated) VALUES(?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, content, status)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PostModel) Get(id int) (*models.Post, error) {
	p := &models.Post{}

	stmt := `SELECT id, title, content, status, created, updated FROM posts WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.Title, &p.Content, &p.Status, &p.Created, &p.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (m *PostModel) Update(id int, title, content, status string) error {
	stmt := `UPDATE posts SET title = ?, content = ?, status = ?, updated = UTC_TIMESTAMP() WHERE id = ?`

	_, err := m.DB.Exec(stmt, title, content, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) Delete(id int) error {
	stmt := `DELETE FROM posts WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) Latest() ([]*models.Post, error) {
	q := `SELECT id, title, content, created, updated FROM posts ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		s := &models.Post{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Updated)
		if err != nil {
			return nil, err
		}

		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
