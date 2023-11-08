package store

import "github.com/artemiyKew/todo-list-rest-api/internal/model"

type UserRepository interface {
	Create(model.User) error
	FindByID(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type WorkRepository interface {
	Create(model.Work) error
	Get(int) ([]*model.Work, error)
	Delete(int, int) error
	Update(model.Work, int) error
}
