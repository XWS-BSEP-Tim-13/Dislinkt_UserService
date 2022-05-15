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
		Biography:   "biography sample",
		IsPrivate:   false,
		Username:    "suki",
		Experiences: []domain.Experience{},
		Educations:  []domain.Education{},
		Skills:      []string{"s1", "s2"},
		Interests:   []primitive.ObjectID{},
		Connections: []primitive.ObjectID{getObjectId("723b0cc3a34d25d8567f9f83"), getObjectId("723b0cc3a34d25d8567f9f84")},
	},
	{
		Id:          getObjectId("723b0cc3a34d25d8567f9f83"),
		FirstName:   "Stefan",
		LastName:    "Ljubovic",
		Email:       "ljubovicstefan@gmail.com",
		PhoneNumber: "0654324995",
		Username:    "ljubo",
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
