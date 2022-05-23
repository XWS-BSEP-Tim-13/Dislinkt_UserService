package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStore interface {
	GetActiveById(id primitive.ObjectID) (*RegisteredUser, error)
	GetAllActive() ([]*RegisteredUser, error)
	Insert(company *RegisteredUser) error
	DeleteAll()
	GetBasicInfo() ([]*RegisteredUser, error)
	Update(user *RegisteredUser) error
	UpdatePersonalInfo(user *RegisteredUser) (primitive.ObjectID, error)
	GetActiveByUsername(username string) (*RegisteredUser, error)
	GetByUsername(username string) (*RegisteredUser, error)
	GetActiveByEmail(email string) (*RegisteredUser, error)
	GetByEmail(email string) (*RegisteredUser, error)
	AddExperience(experience *Experience, userId primitive.ObjectID) error
	AddEducation(education *Education, userId primitive.ObjectID) error
	AddSkill(skill string, userId primitive.ObjectID) error
	RemoveSkill(skill string, userId primitive.ObjectID) error
	AddInterest(companyId primitive.ObjectID, userId primitive.ObjectID) error
	DeleteExperience(experienceId primitive.ObjectID, userId primitive.ObjectID) error
	DeleteEducation(educationId primitive.ObjectID, userId primitive.ObjectID) error
	RemoveInterest(companyId primitive.ObjectID, userId primitive.ObjectID) error
}
