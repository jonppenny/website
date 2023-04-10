package mysql

import (
	"database/sql"
	"errors"

	"jonppenny.co.uk/webapp/pkg/models"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title, content, status, image string) (int, error) {
	stmt := `INSERT INTO posts (title, content, status, created, updated, image) VALUES(?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, title, content, status, image)
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

	stmt := `SELECT id, title, content, status, created, updated, image, excerpt FROM posts WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.Title, &p.Content, &p.Status, &p.Created, &p.Updated, &p.Image, &p.Excerpt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (m *PostModel) Update(id int, title, content, status, image, excerpt string) error {
	stmt := `UPDATE posts SET title = ?, content = ?, status = ?, updated = UTC_TIMESTAMP(), image = ?, excerpt = ? WHERE id = ?`

	_, err := m.DB.Exec(stmt, title, content, status, image, excerpt, id)
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

func (m *PostModel) Latest(limit, offset int) ([]*models.Post, error) {
	rows, err := m.DB.Query("SELECT id, title, content, created, updated, image, excerpt FROM posts ORDER BY created DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		s := &models.Post{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Updated, &s.Image, &s.Excerpt)
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

func (m *PostModel) Total() (int, error) {
	var c int
	err := m.DB.QueryRow("SELECT COUNT(id) FROM posts").Scan(&c)
	if err != nil {
		return 0, err
	}

	return c, nil
}
