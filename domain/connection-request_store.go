package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionRequestStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*ConnectionRequest, error)
	GetAll(ctx context.Context) ([]*ConnectionRequest, error)
	Insert(ctx context.Context, company *ConnectionRequest) error
	DeleteAll(ctx context.Context)
	GetRequestsForUser(ctx context.Context, id primitive.ObjectID) ([]*ConnectionRequest, error)
	Delete(ctx context.Context, id primitive.ObjectID)
}
