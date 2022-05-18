package startup

import (
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var users = []*domain.RegisteredUser{
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f82"),
		FirstName:   "Srdjan",
		LastName:    "Sukovic",
		Email:       "srdjansukovic@gmail.com",
		PhoneNumber: "0649459562",
		Gender:      0,
		DateOfBirth: time.Time{},
		Biography:   "Entrepreneur, investor, and business magnate. Next I’m buying Coca-Cola to put the cocaine back in.",
		IsPrivate:   false,
		Username:    "srdjansukovic",
		Experiences: []domain.Experience{
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d72"),
				Description:        "Senior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Los Angeles",
				CompanyName:        "Google",
				EmploymentType:     0,
				Title:              "Full stack engineer",
			},
			{
				Id:                 getObjectId("723b0cc3a34d25d8567f9d77"),
				Description:        "Junior web engineer in charge of automotive project",
				StartDate:          time.Time{},
				EndDate:            time.Time{},
				Industry:           "Software",
				IsCurrentlyWorking: false,
				Location:           "Las Vegas",
				CompanyName:        "Facebook",
				EmploymentType:     0,
				Title:              "Full stack engineer",
			},
		},
		Educations: []domain.Education{
			{
				Id:           getObjectId("723b0cc3a34d25d8567f9d74"),
				StartDate:    time.Time{},
				EndDate:      time.Time{},
				Description:  "Graduated first in class",
				FieldOfStudy: "Computer science",
				School:       "Harvard",
				Degree:       1,
			},
		},
		Skills:      []string{"s1", "s2"},
		Interests:   []primitive.ObjectID{getObjectId("623b0cc3a34d25d8567f9f82")},
		Connections: []primitive.ObjectID{getObjectId("723b0cc3a34d25d8567f9f83"), getObjectId("723b0cc3a34d25d8567f9f84")},
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f83"),
		FirstName:   "Stefan",
		LastName:    "Ljubovic",
		Email:       "ljubovicstefan@gmail.com",
		PhoneNumber: "0654324995",
		Username:    "stefanljubovic",
		Gender:      0,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   true,
		Experiences: []domain.Experience{},
		Educations:  []domain.Education{},
		Skills:      []string{"s1", "s2"},
		Interests:   []primitive.ObjectID{},
		Connections: []primitive.ObjectID{getObjectId("723b0cc3a34d25d8567f9f84"), getObjectId("723b0cc3a34d25d8567f9f82"), getObjectId("723b0cc3a34d25d8567f9f85")},
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f84"),
		FirstName:   "Ana",
		LastName:    "Gavrilovic",
		Email:       "anagavrilovic@gmail.com",
		PhoneNumber: "0642152",
		Username:    "anagavrilovic",
		Gender:      1,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   false,
		Experiences: []domain.Experience{},
		Educations:  []domain.Education{},
		Skills:      []string{"s1", "s2"},
		Interests:   []primitive.ObjectID{},
		Connections: []primitive.ObjectID{getObjectId("723b0cc3a34d25d8567f9f83"), getObjectId("723b0cc3a34d25d8567f9f82"), getObjectId("723b0cc3a34d25d8567f9f85")},
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f85"),
		FirstName:   "Marija",
		LastName:    "Kljestan",
		Email:       "marijakljestan@gmail.com",
		PhoneNumber: "0642152643",
		Username:    "marijakljestan",
		Gender:      1,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   false,
		Experiences: []domain.Experience{},
		Educations:  []domain.Education{},
		Skills:      []string{"s1", "s2"},
		Interests:   []primitive.ObjectID{},
		Connections: []primitive.ObjectID{getObjectId("723b0cc3a34d25d8567f9f83"), getObjectId("723b0cc3a34d25d8567f9f82")},
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f86"),
		FirstName:   "Lenka",
		LastName:    "Aleksic",
		Email:       "lenka@gmail.com",
		PhoneNumber: "064364364",
		Username:    "lenka",
		Gender:      1,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   false,
		Experiences: []domain.Experience{},
		Educations:  []domain.Education{},
		Skills:      []string{"s1", "s2"},
		Interests:   []primitive.ObjectID{},
		Connections: []primitive.ObjectID{},
	},
}

var connections = []*domain.ConnectionRequest{}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
