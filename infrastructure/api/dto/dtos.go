package dto

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserInfo struct {
	Id          primitive.ObjectID
	FirstName   string
	LastName    string
	Gender      enum.Gender
	DateOfBirth time.Time
	Email       string
	PhoneNumber string
	Biography   string
}

func NewUserInfo(id primitive.ObjectID, fname string, lname string, gender enum.Gender, dob time.Time,
	email string, phoneNumber string, bio string) *UserInfo {

	userInfo := new(UserInfo)
	userInfo.Email = email
	userInfo.Id = id
	userInfo.FirstName = fname
	userInfo.LastName = lname
	userInfo.Gender = gender
	userInfo.Biography = bio
	userInfo.DateOfBirth = dob
	userInfo.PhoneNumber = phoneNumber

	return userInfo
}

func NewUserFromUserInfo(userInfo UserInfo) *domain.RegisteredUser {
	user := new(domain.RegisteredUser)
	user.Id = userInfo.Id
	user.LastName = userInfo.LastName
	user.Gender = userInfo.Gender
	user.DateOfBirth = userInfo.DateOfBirth
	user.Email = userInfo.Email
	user.PhoneNumber = userInfo.PhoneNumber
	user.Biography = userInfo.Biography

	return user
}
