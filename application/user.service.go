package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	logger "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/logging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/tracer"
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

func (service *UserService) Get(ctx context.Context, id primitive.ObjectID) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetActiveById(ctx, id)
}

func (service *UserService) DeleteConnectionRequest(ctx context.Context, connectionId primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE DeleteConnectionRequest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	service.connectionStore.Delete(ctx, connectionId)
}

func (service *UserService) GetRequestsForUser(ctx context.Context, id primitive.ObjectID) ([]*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE GetRequestsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	resp, err := service.connectionStore.GetRequestsForUser(ctx, id)
	return resp, err
}

func (service *UserService) FindByFilter(ctx context.Context, filter string) ([]*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE FindByFilter")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	users, err := service.store.GetAllActive(ctx)
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

func (service *UserService) GetAll(ctx context.Context) ([]*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE GetAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllActive(ctx)
}

func (service *UserService) UpdatePersonalInfo(ctx context.Context, user *domain.RegisteredUser) (primitive.ObjectID, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE UpdatePersonalInfo")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.UpdatePersonalInfo(ctx, user)
}

func (service *UserService) CreateNewUser(ctx context.Context, user *domain.RegisteredUser) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE CreateNewUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	dbUser, _ := service.store.GetByUsername(ctx, (*user).Username)
	if dbUser != nil {
		err := errors.New("username already exists")
		return nil, err
	}

	dbUser, _ = service.store.GetByEmail(ctx, (*user).Email)
	if dbUser != nil {
		err := errors.New("email already exists")
		return nil, err
	}

	(*user).Id = primitive.NewObjectID()
	(*user).IsActive = false
	err := service.store.Insert(ctx, user)
	if err != nil {
		err := errors.New("error while creating new user")
		return nil, err
	}

	return user, nil
}

func (service *UserService) AddExperience(ctx context.Context, experience *domain.Experience, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE AddExperience")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	experience.Id = primitive.NewObjectID()
	return service.store.AddExperience(ctx, experience, userId)
}

func (service *UserService) AddEducation(ctx context.Context, education *domain.Education, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE AddEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	education.Id = primitive.NewObjectID()
	return service.store.AddEducation(ctx, education, userId)
}

func (service *UserService) AddSkill(ctx context.Context, skill string, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE AddSkill")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.AddSkill(ctx, skill, userId)
}

func (service *UserService) RemoveSkill(ctx context.Context, skill string, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE RemoveSkill")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.RemoveSkill(ctx, skill, userId)
}

func (service *UserService) AddInterest(ctx context.Context, companyId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE AddInterest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.AddInterest(ctx, companyId, userId)
}

func (service *UserService) DeleteExperience(ctx context.Context, experienceId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE DeleteExperience")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.DeleteExperience(ctx, experienceId, userId)
}

func (service *UserService) DeleteEducation(ctx context.Context, educationId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE DeleteEducation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.DeleteEducation(ctx, educationId, userId)
}
func (service *UserService) RemoveInterest(ctx context.Context, companyId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE RemoveInterest")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.RemoveInterest(ctx, companyId, userId)
}
func (service *UserService) GetByUsername(ctx context.Context, username string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE GetByUsername")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := service.store.GetActiveByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) GetByEmail(ctx context.Context, email string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE GetByEmail")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := service.store.GetActiveByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) ActivateAccount(ctx context.Context, email string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE ActivateAccount")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := service.store.UpdateIsActive(ctx, email)
	if err != nil {
		err := errors.New("error activating account")
		return "", err
	}

	return "Account successfully activated!", nil
}

func (service *UserService) ChangeAccountPrivacy(ctx context.Context, username string, isPrivate bool) error {
	span := tracer.StartSpanFromContext(ctx, "SERVICE ChangeAccountPrivacy")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.ChangeAccountPrivacy(ctx, isPrivate, username)
}

func (service *UserService) GetNotificationsForUser(ctx context.Context, username string) ([]*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "SERVICE ChangeAccountPrivacy")
	defer span.Finish()
	user, err := service.store.GetByUsername(ctx, username)
	fmt.Println(user)
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

func (service *UserService) UpdateNotificationAlert(username string, displayNotifications bool) error {
	return service.store.UpdateDisplayUserNotifications(displayNotifications, username)
}
