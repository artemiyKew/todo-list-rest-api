package sqlstore

import (
	"database/sql"
	"errors"

	"github.com/artemiyKew/todo-list-rest-api/internal/model"
)

type WorkRepository struct {
	store *Store
}

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}
	return r.store.db.
		QueryRow(
			"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
			u.Email,
			u.EncryptedPassword).
		Scan(&u.ID)
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.
		QueryRow("SELECT id, email, encrypted_password FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("err no rows")
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.
		QueryRow("SELECT id, email, encrypted_password FROM users WHERE email = $1", email).
		Scan(&u.ID, &u.Email, &u.EncryptedPassword); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *WorkRepository) Create(w *model.Work) error {
	return nil
}

func (r *WorkRepository) Change(w *model.Work) error {
	return nil
}

func (r *WorkRepository) Delete(id int) error {
	return nil
}
