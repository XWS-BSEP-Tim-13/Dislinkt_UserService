package persistence

import (
	"context"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/tracer"
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

func (store ConnectionsMongoDBStore) Delete(ctx context.Context, id primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "DB Delete")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"_id": id}
	store.connections.DeleteOne(context.TODO(), filter)
}

func (store ConnectionsMongoDBStore) GetRequestsForUser(ctx context.Context, id primitive.ObjectID) ([]*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetRequestsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{"to._id", id}}
	return store.filter(ctx, filter)
}

func (store ConnectionsMongoDBStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"_id": id}
	return store.filterOne(ctx, filter)
}

func (store ConnectionsMongoDBStore) GetAll(ctx context.Context) ([]*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB GetAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}}
	return store.filter(ctx, filter)
}

func (store ConnectionsMongoDBStore) Insert(ctx context.Context, connection *domain.ConnectionRequest) error {
	span := tracer.StartSpanFromContext(ctx, "DB Insert")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.connections.InsertOne(context.TODO(), connection)
	if err != nil {
		return err
	}
	connection.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store ConnectionsMongoDBStore) DeleteAll(ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "DB DeleteAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	store.connections.DeleteMany(ctx, bson.D{{}})
}

func (store *ConnectionsMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*domain.ConnectionRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "DB filter")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	cursor, err := store.connections.Find(ctx, filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decodeConnection(cursor)
}

func (store *ConnectionsMongoDBStore) filterOne(ctx context.Context, filter interface{}) (connection *domain.ConnectionRequest, err error) {
	span := tracer.StartSpanFromContext(ctx, "DB filterOne")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result := store.connections.FindOne(ctx, filter)
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
