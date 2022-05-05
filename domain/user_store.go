package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStore interface {
	Get(id primitive.ObjectID) (*RegisteredUser, error)
	GetAll() ([]*RegisteredUser, error)
	Insert(company *RegisteredUser) error
	DeleteAll()
}
