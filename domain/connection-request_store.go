package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConnectionRequestStore interface {
	Get(id primitive.ObjectID) (*ConnectionRequest, error)
	GetAll() ([]*ConnectionRequest, error)
	Insert(company *ConnectionRequest) error
	DeleteAll()
	GetRequestsForUser(id primitive.ObjectID) ([]*ConnectionRequest, error)
	Delete(id primitive.ObjectID)
}
