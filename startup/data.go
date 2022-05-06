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
		IsPrivate:   true,
		Experiences: []domain.Experience{},
		Educations:  []domain.Education{},
		Skills:      []string{"s1", "s2"},
		Interests:   []primitive.ObjectID{},
		Connections: []primitive.ObjectID{},
	},
	{
		Id:          getObjectId("f880e6be-cd34-11ec-9d64-0242ac120002"),
		FirstName:   "Stefan",
		LastName:    "Ljubovic",
		Email:       "ljubovicstefan@gmail.com",
		PhoneNumber: "0654324995",
		Gender:      0,
		DateOfBirth: time.Time{},
		Biography:   "biography sample",
		IsPrivate:   true,
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
