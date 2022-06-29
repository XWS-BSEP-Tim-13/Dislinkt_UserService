package domain

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Education struct {
	Id           primitive.ObjectID `bson:"_id"`
	School       string             `bson:"school" validate:"required,alphaNumSigns"`
	Degree       enum.Degree        `bson:"degree"`
	FieldOfStudy string             `bson:"field_of_study" validate:"required,alphaNumSigns"`
	StartDate    time.Time          `bson:"start_date" validate:"required"`
	EndDate      time.Time          `bson:"end_date" validate:"required"`
	Description  string             `bson:"description"`
}

type Experience struct {
	Id                 primitive.ObjectID  `bson:"_id"`
	Title              string              `bson:"title" validate:"required,alphaNumSigns"`
	EmploymentType     enum.EmploymentType `bson:"employment_type"`
	CompanyName        string              `bson:"company_name" validate:"required,alphaNumSigns"`
	Location           string              `bson:"location" validate:"alphaNumSigns"`
	IsCurrentlyWorking bool                `bson:"is_currently_working"`
	StartDate          time.Time           `bson:"start_date" validate:"required"`
	EndDate            time.Time           `bson:"end_date" validate:"required"`
	Industry           string              `bson:"industry" validate:"alphaNumSigns"`
	Description        string              `bson:"description"`
}

type RegisteredUser struct {
	Id           primitive.ObjectID   `bson:"_id"`
	FirstName    string               `bson:"first_name" validate:"required,alphaSpace"`
	LastName     string               `bson:"last_name" validate:"required,alphaSpace"`
	Email        string               `bson:"email" validate:"required,email"`
	PhoneNumber  string               `bson:"phone_number" validate:"required,numeric,min=9,max=10"`
	Gender       enum.Gender          `bson:"gender"`
	DateOfBirth  time.Time            `bson:"date_of_birth" validate:"required"`
	Biography    string               `bson:"biography" validate:"required,max=256"`
	IsPrivate    bool                 `bson:"is_private"`
	IsActive     bool                 `bson:"is_active"`
	Experiences  []Experience         `bson:"experiences"`
	Educations   []Education          `bson:"educations"`
	Skills       []string             `bson:"skills"`
	Interests    []primitive.ObjectID `bson:"interests"`
	Connections  []primitive.ObjectID `bson:"connections"`
	Username     string               `bson:"username" validate:"required,username"`
	BlockedUsers []string             `bson:"blocked_users"`
}

type ConnectionRequest struct {
	Id          primitive.ObjectID `bson:"_id"`
	From        RegisteredUser     `bson:"from"`
	To          RegisteredUser     `bson:"to"`
	RequestTime time.Time          `bson:"request_time"`
}
