package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	GetActiveById(ctx context.Context, id primitive.ObjectID) (*RegisteredUser, error)
	GetAllActive(ctx context.Context) ([]*RegisteredUser, error)
	Insert(ctx context.Context, company *RegisteredUser) error
	DeleteAll(ctx context.Context)
	GetBasicInfo(ctx context.Context) ([]*RegisteredUser, error)
	Update(ctx context.Context, user *RegisteredUser) error
	UpdateIsActive(ctx context.Context, email string) error
	UpdatePersonalInfo(ctx context.Context, user *RegisteredUser) (primitive.ObjectID, error)
	GetActiveByUsername(ctx context.Context, username string) (*RegisteredUser, error)
	GetByUsername(ctx context.Context, username string) (*RegisteredUser, error)
	GetActiveByEmail(ctx context.Context, email string) (*RegisteredUser, error)
	GetByEmail(ctx context.Context, email string) (*RegisteredUser, error)
	AddExperience(ctx context.Context, experience *Experience, userId primitive.ObjectID) error
	AddEducation(ctx context.Context, education *Education, userId primitive.ObjectID) error
	AddSkill(ctx context.Context, skill string, userId primitive.ObjectID) error
	RemoveSkill(ctx context.Context, skill string, userId primitive.ObjectID) error
	AddInterest(ctx context.Context, companyId primitive.ObjectID, userId primitive.ObjectID) error
	DeleteExperience(ctx context.Context, experienceId primitive.ObjectID, userId primitive.ObjectID) error
	DeleteEducation(ctx context.Context, educationId primitive.ObjectID, userId primitive.ObjectID) error
	RemoveInterest(ctx context.Context, companyId primitive.ObjectID, userId primitive.ObjectID) error
	ChangeAccountPrivacy(ctx context.Context, isPrivate bool, username string) error
	UpdateDisplayUserNotifications(displayNotifications bool, username string) error
}
