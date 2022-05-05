package application

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Get(id primitive.ObjectID) (*domain.RegisteredUser, error) {
	return service.store.Get(id)
}

func (service *UserService) FindByFilter(filter string) ([]*domain.RegisteredUser, error) {
	users, err := service.store.GetAll()
	if err != nil {
		return nil, err
	}
	var filteredUsers []*domain.RegisteredUser
	for _, user := range users {
		fullName := user.FirstName + " " + user.LastName
		if strings.Contains(strings.ToLower(fullName), strings.ToLower(filter)) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, nil
}

func (service *UserService) GetAll() ([]*domain.RegisteredUser, error) {
	return service.store.GetAll()
}
