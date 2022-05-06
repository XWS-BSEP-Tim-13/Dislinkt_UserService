package application

import (
	"errors"
	"fmt"
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
	fromUser, _ := service.store.Get(idFrom)
	if err != nil {
		return err
	}
	fmt.Printf("In service trace: \n")
	if toUser.IsPrivate {
		var request = domain.ConnectionRequest{
			From:        *fromUser,
			To:          *toUser,
			RequestTime: time.Now(),
		}
		service.connectionStore.Insert(&request)
	} else {
		toUser.Connections = append(toUser.Connections, idFrom)
		service.store.Update(toUser)
	}
	fmt.Printf("Saved to db: \n")
	return nil
}

func (service *UserService) GetRequestsForUser(id primitive.ObjectID) ([]*domain.ConnectionRequest, error) {
	return service.connectionStore.GetRequestsForUser(id)
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

func (service *UserService) UpdatePersonalInfo(user *domain.RegisteredUser) (primitive.ObjectID, error) {
	return service.store.UpdatePersonalInfo(user)
}

func (service *UserService) CreateNewUser(user *domain.RegisteredUser) (*domain.RegisteredUser, error) {
	//dbUser, _ := service.store.GetByUsername((*user).Username)
	//if (*dbUser).Username != "" {
	//	err := errors.New("username already exists")
	//	return nil, err
	//}
	(*user).Id = primitive.NewObjectID()
	err := service.store.Insert(user)
	if err != nil {
		err := errors.New("error while creating new user")
		return nil, err
	}

	return user, nil
}
