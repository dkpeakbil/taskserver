package mem

import (
	"encoding/json"
	"github.com/dkpeakbil/taskserver/domain"
	"github.com/dkpeakbil/taskserver/repository"
	"sync"
	"time"
)

type inMemRepository struct {
	sync.RWMutex
	seq   int
	users map[int]*domain.User
}

func NewInMemoryRepository() (repository.Repository, error) {
	return &inMemRepository{
		seq:   1,
		users: make(map[int]*domain.User),
	}, nil
}

func (i *inMemRepository) Save(user *domain.User) (*domain.User, error) {
	i.Lock()
	defer i.Unlock()

	user.ID = i.seq
	_, ok := i.users[user.ID]
	if ok {
		return nil, domain.ErrUserAlreadyExists
	}

	for _, u := range i.users {
		if u.Username == user.Username {
			return nil, domain.ErrUsernameHasTaken
		}
	}

	now := time.Now()
	user.UpdatedAt = now
	user.CreatedAt = now

	i.users[user.ID] = user

	i.seq++

	return user, nil
}

func (i *inMemRepository) FindByID(id int) (*domain.User, error) {
	i.RLock()
	defer i.RUnlock()

	for _, u := range i.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (i *inMemRepository) FindByUsername(username string) (*domain.User, error) {
	i.RLock()
	defer i.RUnlock()

	for _, u := range i.users {
		if u.Username == username {
			return u, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (i *inMemRepository) String() string {
	i.RLock()
	defer i.RUnlock()

	j, _ := json.Marshal(i.users)
	return string(j)
}
