package application

import (
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

func (service *UserService) CheckIfUserCanReadPosts(idFrom, idTo primitive.ObjectID) (bool, error) {
	toUser, err := service.store.Get(idTo)
	if err != nil {
		return false, err
	}
	if !toUser.IsPrivate {
		return true, nil
	}
	for _, conId := range toUser.Connections {
		if conId == idFrom {
			return true, nil
		}
	}
	return false, nil
}

func (service *UserService) AcceptConnection(connectionId primitive.ObjectID) error {
	connection, err := service.connectionStore.Get(connectionId)
	if err != nil {
		return err
	}
	connection.To.Connections = append(connection.To.Connections, connection.From.Id)
	fmt.Printf("Saved connection %s \n", connection.To.Connections)
	err1 := service.store.Update(&connection.To)
	if err != nil {
		return err1
	}
	service.connectionStore.Delete(connectionId)
	return nil
}

func (service *UserService) GetRequestsForUser(id primitive.ObjectID) ([]*domain.ConnectionRequest, error) {
	resp, err := service.connectionStore.GetRequestsForUser(id)
	fmt.Printf("Response %d\n", len(resp))
	return resp, err
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
