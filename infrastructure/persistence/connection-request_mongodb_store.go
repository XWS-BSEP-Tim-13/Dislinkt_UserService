package persistence

import (
	"context"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CON_DATABASE   = "connections"
	CON_COLLECTION = "connection"
)

type ConnectionsMongoDBStore struct {
	connections *mongo.Collection
}

func (store ConnectionsMongoDBStore) GetRequestsForUser(id primitive.ObjectID) ([]*domain.ConnectionRequest, error) {
	filter := bson.D{{"to._id", id}}
	return store.filter(filter)
}

func (store ConnectionsMongoDBStore) Get(id primitive.ObjectID) (*domain.ConnectionRequest, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store ConnectionsMongoDBStore) GetAll() ([]*domain.ConnectionRequest, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store ConnectionsMongoDBStore) Insert(connection *domain.ConnectionRequest) error {
	result, err := store.connections.InsertOne(context.TODO(), connection)
	if err != nil {
		return err
	}
	connection.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store ConnectionsMongoDBStore) DeleteAll() {
	store.connections.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *ConnectionsMongoDBStore) filter(filter interface{}) ([]*domain.ConnectionRequest, error) {
	cursor, err := store.connections.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decodeConnection(cursor)
}

func (store *ConnectionsMongoDBStore) filterOne(filter interface{}) (connection *domain.ConnectionRequest, err error) {
	result := store.connections.FindOne(context.TODO(), filter)
	err = result.Decode(&connection)
	return
}

func decodeConnection(cursor *mongo.Cursor) (connections []*domain.ConnectionRequest, err error) {
	for cursor.Next(context.TODO()) {
		var connection domain.ConnectionRequest
		err = cursor.Decode(&connection)
		if err != nil {
			return
		}
		connections = append(connections, &connection)
	}
	err = cursor.Err()
	return
}

func NewConnectionsMongoDBStore(client *mongo.Client) domain.ConnectionRequestStore {
	connections := client.Database(CON_DATABASE).Collection(CON_COLLECTION)
	return &ConnectionsMongoDBStore{
		connections: connections,
	}
}
