package application

import (
	"errors"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type UserService struct {
	store             domain.UserStore
	connectionStore   domain.ConnectionRequestStore
	logger            *logger.Logger
	notificationStore domain.NotificationStore
}

func NewUserService(store domain.UserStore, connectionStore domain.ConnectionRequestStore, logger *logger.Logger, notificationStore domain.NotificationStore) *UserService {
	return &UserService{
		store:             store,
		connectionStore:   connectionStore,
		logger:            logger,
		notificationStore: notificationStore,
	}
}

func (service *UserService) Get(id primitive.ObjectID) (*domain.RegisteredUser, error) {
	return service.store.GetActiveById(id)
}

//func (service *UserService) GetConnectionUsernamesForUser(username string) ([]string, error) {
//	user, err := service.store.GetActiveByUsername(username)
//	if err != nil {
//		fmt.Println("Active error")
//		return nil, err
//	}
//	var retVal []string
//	for _, conId := range user.Connections {
//		conUser, _ := service.store.GetActiveById(conId)
//		retVal = append(retVal, conUser.Username)
//		fmt.Printf("Username : %s\n", conUser.Username)
//	}
//	retVal = append(retVal, username)
//	return retVal, nil
//}

func (service *UserService) DeleteConnectionRequest(connectionId primitive.ObjectID) {
	service.connectionStore.Delete(connectionId)
}

func (service *UserService) GetRequestsForUser(id primitive.ObjectID) ([]*domain.ConnectionRequest, error) {
	resp, err := service.connectionStore.GetRequestsForUser(id)
	return resp, err
}

func (service *UserService) FindByFilter(filter string) ([]*domain.RegisteredUser, error) {
	users, err := service.store.GetAllActive()
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
	return service.store.GetAllActive()
}

func (service *UserService) UpdatePersonalInfo(user *domain.RegisteredUser) (primitive.ObjectID, error) {
	return service.store.UpdatePersonalInfo(user)
}

func (service *UserService) CreateNewUser(user *domain.RegisteredUser) (*domain.RegisteredUser, error) {
	dbUser, _ := service.store.GetByUsername((*user).Username)
	if dbUser != nil {
		err := errors.New("username already exists")
		return nil, err
	}

	dbUser, _ = service.store.GetByEmail((*user).Email)
	if dbUser != nil {
		err := errors.New("email already exists")
		return nil, err
	}

	(*user).Id = primitive.NewObjectID()
	(*user).IsActive = false
	err := service.store.Insert(user)
	if err != nil {
		err := errors.New("error while creating new user")
		return nil, err
	}

	return user, nil
}

func (service *UserService) AddExperience(experience *domain.Experience, userId primitive.ObjectID) error {
	experience.Id = primitive.NewObjectID()
	return service.store.AddExperience(experience, userId)
}

func (service *UserService) AddEducation(education *domain.Education, userId primitive.ObjectID) error {
	education.Id = primitive.NewObjectID()
	return service.store.AddEducation(education, userId)
}

func (service *UserService) AddSkill(skill string, userId primitive.ObjectID) error {
	return service.store.AddSkill(skill, userId)
}

func (service *UserService) RemoveSkill(skill string, userId primitive.ObjectID) error {
	return service.store.RemoveSkill(skill, userId)
}

func (service *UserService) AddInterest(companyId primitive.ObjectID, userId primitive.ObjectID) error {
	return service.store.AddInterest(companyId, userId)
}

func (service *UserService) DeleteExperience(experienceId primitive.ObjectID, userId primitive.ObjectID) error {
	return service.store.DeleteExperience(experienceId, userId)
}

func (service *UserService) DeleteEducation(educationId primitive.ObjectID, userId primitive.ObjectID) error {
	return service.store.DeleteEducation(educationId, userId)
}
func (service *UserService) RemoveInterest(companyId primitive.ObjectID, userId primitive.ObjectID) error {
	return service.store.RemoveInterest(companyId, userId)
}
func (service *UserService) GetByUsername(username string) (*domain.RegisteredUser, error) {
	user, err := service.store.GetActiveByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) GetByEmail(email string) (*domain.RegisteredUser, error) {
	user, err := service.store.GetActiveByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) ActivateAccount(email string) (string, error) {
	err := service.store.UpdateIsActive(email)
	if err != nil {
		err := errors.New("error activating account")
		return "", err
	}

	return "Account successfully activated!", nil
}

func (service *UserService) ChangeAccountPrivacy(username string, isPrivate bool) error {
	return service.store.ChangeAccountPrivacy(isPrivate, username)
}

func (service *UserService) GetNotificationsForUser(username string) ([]*domain.Notification, error) {
	user, err := service.store.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	var notifications []*domain.Notification
	for _, usernameFrom := range user.Connections {
		fmt.Println(usernameFrom)
		notificationsDb, _ := service.notificationStore.GetByUsername(usernameFrom)
		notifications = append(notifications, notificationsDb...)
	}
	return notifications, nil
}

func (service *UserService) SaveNotification(notification *domain.Notification) error {
	err := service.notificationStore.Insert(notification)
	if err != nil {
		return err
	}
	return nil
}
