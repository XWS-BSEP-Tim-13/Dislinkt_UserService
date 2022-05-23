package domain

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Education struct {
	Id           primitive.ObjectID `bson:"_id"`
	School       string             `bson:"school"`
	Degree       enum.Degree        `bson:"degree"`
	FieldOfStudy string             `bson:"field_of_study"`
	StartDate    time.Time          `bson:"start_date"`
	EndDate      time.Time          `bson:"end_date"`
	Description  string             `bson:"description"`
}

type Experience struct {
	Id                 primitive.ObjectID  `bson:"_id"`
	Title              string              `bson:"title"`
	EmploymentType     enum.EmploymentType `bson:"employment_type"`
	CompanyName        string              `bson:"company_name"`
	Location           string              `bson:"location"`
	IsCurrentlyWorking bool                `bson:"is_currently_working"`
	StartDate          time.Time           `bson:"start_date"`
	EndDate            time.Time           `bson:"end_date"`
	Industry           string              `bson:"industry"`
	Description        string              `bson:"description"`
}

type RegisteredUser struct {
	Id          primitive.ObjectID   `bson:"_id"`
	FirstName   string               `bson:"first_name"`
	LastName    string               `bson:"last_name"`
	Email       string               `bson:"email"`
	PhoneNumber string               `bson:"phone_number"`
	Gender      enum.Gender          `bson:"gender"`
	DateOfBirth time.Time            `bson:"date_of_birth"`
	Biography   string               `bson:"biography"`
	IsPrivate   bool                 `bson:"is_private"`
	IsActive    bool                 `bson:"is_active"`
	Experiences []Experience         `bson:"experiences"`
	Educations  []Education          `bson:"educations"`
	Skills      []string             `bson:"skills"`
	Interests   []primitive.ObjectID `bson:"interests"`
	Connections []primitive.ObjectID `bson:"connections"`
	Username    string               `bson:"username"`
}

type ConnectionRequest struct {
	Id          primitive.ObjectID `bson:"_id"`
	From        RegisteredUser     `bson:"from"`
	To          RegisteredUser     `bson:"to"`
	RequestTime time.Time          `bson:"request_time"`
}

func NewRegisteredUser(id primitive.ObjectID, fname string, lname string, gender enum.Gender, dob time.Time, email string,
	phoneNumber string, biography string) {

}
