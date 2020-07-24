package usecase

import (
	"github.com/dkpeakbil/taskserver/domain"
)

type UseCase interface {
	Register(*domain.RegisterRequest) *domain.RegisterResponse
	Auth(request *domain.AuthRequest) *domain.AuthResponse
	AuthGame(request *domain.AuthGameRequest) *domain.AuthGameResponse
	GetUsers(request *domain.GetUsersRequest) *domain.GetUsersResponse
}
