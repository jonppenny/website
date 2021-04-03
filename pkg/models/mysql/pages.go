package mysql

import (
	"database/sql"
	"errors"
	"jonppenny.co.uk/webapp/pkg/models"
)

type PageModel struct {
	DB *sql.DB
}

func (m *PageModel) Insert(title, content, status, slug string) (int, error) {
	stmt := `INSERT INTO pages (title, content, status, slug, created, updated) VALUES (?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, content, status, slug)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PageModel) Get(id int) (*models.Page, error) {
	p := &models.Page{}

	stmt := `SELECT id, title, content, status, slug, created, updated FROM pages WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.Title, &p.Content, &p.Status, &p.Slug, &p.Created, &p.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (m *PageModel) GetBySlug(slug string) (*models.Page, error) {
	p := &models.Page{}

	stmt := `SELECT id, title, content, status, slug, created, updated FROM pages WHERE slug = ?`
	err := m.DB.QueryRow(stmt, slug).Scan(&p.ID, &p.Title, &p.Content, &p.Status, &p.Slug, &p.Created, &p.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return p, nil
}

func (m *PageModel) Update(id int, title, content, status, slug string) error {
	stmt := `UPDATE pages SET title = ?, content = ?, status = ?, slug = ?, updated = UTC_TIMESTAMP() WHERE id = ?`

	_, err := m.DB.Exec(stmt, title, content, status, slug, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PageModel) Delete(id int) error {
	stmt := `DELETE FROM pages WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *PageModel) GetAll() ([]*models.Page, error) {
	q := `SELECT id, title, content, status, slug, created, updated FROM pages ORDER BY id`

	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pages []*models.Page

	for rows.Next() {
		s := &models.Page{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Status, &s.Slug, &s.Created, &s.Updated)
		if err != nil {
			return nil, err
		}

		pages = append(pages, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}
