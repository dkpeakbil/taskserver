package impl

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dkpeakbil/taskserver/domain"
	"github.com/dkpeakbil/taskserver/repository"
	"github.com/dkpeakbil/taskserver/usecase"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"syreclabs.com/go/faker"
	"time"
)

type useCase struct {
	repo       repository.Repository
	validation *validator.Validate
}

func NewUseCase(repo repository.Repository, validation *validator.Validate) (usecase.UseCase, error) {
	return &useCase{
		repo:       repo,
		validation: validation,
	}, nil
}

func (u *useCase) Register(request *domain.RegisterRequest) *domain.RegisterResponse {
	var err error
	if err = u.validation.Struct(request); err != nil {
		return &domain.RegisterResponse{Status: false}
	}

	password := md5.Sum([]byte(request.Password))
	user := &domain.User{
		Username: request.Username,
		Password: hex.EncodeToString(password[:]),
	}

	if user, err = u.repo.Save(user); err != nil {
		return &domain.RegisterResponse{Status: false}
	}

	return &domain.RegisterResponse{Status: true}
}

func (u *useCase) Auth(request *domain.AuthRequest) *domain.AuthResponse {
	var err error
	if err = u.validation.Struct(request); err != nil {
		return &domain.AuthResponse{Status: false}
	}

	user, err := u.repo.FindByUsername(request.Username)
	if err != nil {
		return &domain.AuthResponse{Status: false}
	}

	password := md5.Sum([]byte(request.Password))
	if user.Password != hex.EncodeToString(password[:]) {
		return &domain.AuthResponse{Status: false}
	}

	claims := domain.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(domain.TokenSecret)
	if err != nil {
		return &domain.AuthResponse{Status: false}
	}

	return &domain.AuthResponse{
		Status: true,
		Token:  signed,
	}
}

func (u *useCase) AuthGame(request *domain.AuthGameRequest) *domain.AuthGameResponse {
	claims := &domain.TokenClaims{}
	token, err := jwt.ParseWithClaims(request.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return domain.TokenSecret, nil
	})

	if err != nil {
		logrus.Debugf("err parsing claims %s", err)
		return &domain.AuthGameResponse{
			Status: false,
		}
	}

	if !token.Valid {
		return &domain.AuthGameResponse{
			Status: false,
		}
	}

	user, err := u.repo.FindByID(claims.UserID)
	if err != nil {
		return &domain.AuthGameResponse{Status: false}
	}

	return &domain.AuthGameResponse{
		Status:   true,
		Username: user.Username,
		Message:  fmt.Sprintf("Welcome, %s", user.Username),
	}
}

func (u *useCase) GetUsers(request *domain.GetUsersRequest) *domain.GetUsersResponse {
	if request.Limit < 10 {
		request.Limit = 10
	}

	if request.Offset < 0 {
		request.Offset = 0
	}

	var err error
	if err = u.validation.Struct(request); err != nil {
		return &domain.GetUsersResponse{Status: false}
	}

	response := &domain.GetUsersResponse{
		Status: true,
		Users:  make([]*domain.DummyUser, request.Limit),
		Limit:  request.Limit,
		Offset: request.Offset,
	}

	for i := 0; i < request.Limit; i++ {
		dummyUser := &domain.DummyUser{
			Username: faker.Name().FirstName(),
			Avatar:   faker.Avatar().String(),
		}

		response.Users[i] = dummyUser
	}

	return response
}
