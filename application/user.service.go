package application

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type UserService struct {
	store           domain.UserStore
	connectionStore domain.ConnectionRequestStore
}

func NewUserService(store domain.UserStore, connectionStore domain.ConnectionRequestStore) *UserService {
	return &UserService{
		store:           store,
		connectionStore: connectionStore,
	}
}

func (service *UserService) Get(id primitive.ObjectID) (*domain.RegisteredUser, error) {
	return service.store.Get(id)
}

func (service *UserService) RequestConnection(idFrom, idTo primitive.ObjectID) error {
	toUser, err := service.store.Get(idTo)
	if err != nil {
		return err
	}
	if toUser.IsPrivate {
		fromUser, _ := service.store.Get(idFrom)
		var request = domain.ConnectionRequest{
			From:        *fromUser,
			To:          *toUser,
			RequestTime: time.Now(),
		}
		service.connectionStore.Insert(&request)
	} else {
		toUser.Connections = append(toUser.Connections, idFrom)
	}
	return nil
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
