package sqlstore

import (
	"errors"

	"github.com/artemiyKew/todo-list-rest-api/internal/model"
)

type WorkRepository struct {
	store *Store
}

func (r *WorkRepository) Create(w model.Work) error {
	if err := w.BeforeCreate(); err != nil {
		return err
	}

	return r.store.
		db.
		QueryRow("INSERT INTO works (user_id, name, description, created_at, expired_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			w.User_ID, w.Name, w.Description, w.CreatedAt, w.ExpiredAt).
		Scan(&w.ID)
}

func (r *WorkRepository) Get(userid int) ([]*model.Work, error) {
	w := make([]*model.Work, 0)
	rows, err := r.store.db.Query("SELECT * FROM works WHERE user_id = $1", userid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		work := &model.Work{}
		if err := rows.
			Scan(&work.ID, &work.User_ID, &work.Name,
				&work.Description, &work.CreatedAt, &work.ExpiredAt); err != nil {
			return nil, err
		}
		w = append(w, work)
	}

	return w, nil
}

func (r *WorkRepository) Update(w model.Work, work_id int) error {
	if w.IsEmpty() {
		return errors.New("empty data")
	}

	if work_id < 0 {
		return errors.New("out of range")
	}

	if err := r.store.db.
		QueryRow("UPDATE works SET name = $1, description = $2 WHERE id = $3 AND user_id = $4",
			w.Name, w.Description, work_id, w.User_ID).Err(); err != nil {
		return err
	}

	return nil
}

func (r *WorkRepository) Delete(work_id, user_id int) error {
	if work_id < 0 {
		return errors.New("out of range")
	}

	if err := r.store.db.
		QueryRow("DELETE FROM works WHERE id = $1 AND user_id = $2",
			work_id, user_id).Err(); err != nil {
		return err
	}

	return nil
}
