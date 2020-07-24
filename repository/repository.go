package repository

import (
	"github.com/dkpeakbil/taskserver/domain"
)

type Repository interface {
	Save(user *domain.User) (*domain.User, error)
	FindByID(id int) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
}
