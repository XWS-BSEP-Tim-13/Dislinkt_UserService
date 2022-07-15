package persistence

import (
	"context"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE_NOTIFICATIONS   = "notifications"
	COLLECTION_NOTIFICATIONS = "notification"
)

type NotificationMongoDBStore struct {
	notifications *mongo.Collection
}

func NewNotificationMongoDBStore(client *mongo.Client) domain.NotificationStore {
	notifications := client.Database(DATABASE_NOTIFICATIONS).Collection(COLLECTION_NOTIFICATIONS)
	return &NotificationMongoDBStore{
		notifications: notifications,
	}
}

func (store NotificationMongoDBStore) Insert(notification *domain.Notification) error {
	result, err := store.notifications.InsertOne(context.TODO(), notification)
	if err != nil {
		return err
	}
	notification.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store NotificationMongoDBStore) DeleteAll() {
	store.notifications.DeleteMany(context.TODO(), bson.D{{}})
}

func (store NotificationMongoDBStore) GetByUsername(username string) ([]*domain.Notification, error) {
	filter := bson.M{"username": username}
	return store.filter(filter)
}

func (store *NotificationMongoDBStore) filter(filter interface{}) ([]*domain.Notification, error) {
	cursor, err := store.notifications.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decodeNotifications(cursor)
}

func (store *NotificationMongoDBStore) filterOne(filter interface{}) (notification *domain.Notification, err error) {
	result := store.notifications.FindOne(context.TODO(), filter)
	err = result.Decode(&notification)
	return
}

func decodeNotifications(cursor *mongo.Cursor) (notifications []*domain.Notification, err error) {
	for cursor.Next(context.TODO()) {
		var notification domain.Notification
		err = cursor.Decode(&notification)
		if err != nil {
			return
		}
		notifications = append(notifications, &notification)
	}
	err = cursor.Err()
	return
}
