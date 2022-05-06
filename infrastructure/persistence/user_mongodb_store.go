package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DATABASE   = "users"
	COLLECTION = "user"
)

type UserMongoDBStore struct {
	users *mongo.Collection
}

func NewUserMongoDBStore(client *mongo.Client) domain.UserStore {
	users := client.Database(DATABASE).Collection(COLLECTION)
	return &UserMongoDBStore{
		users: users,
	}
}

func (store *UserMongoDBStore) Get(id primitive.ObjectID) (*domain.RegisteredUser, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetAll() ([]*domain.RegisteredUser, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) GetByUsername(username string) (*domain.RegisteredUser, error) {
	filter := bson.M{"username": username}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) Insert(user *domain.RegisteredUser) error {
	result, err := store.users.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *UserMongoDBStore) DeleteAll() {
	store.users.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *UserMongoDBStore) filter(filter interface{}) ([]*domain.RegisteredUser, error) {
	cursor, err := store.users.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *UserMongoDBStore) Update(user *domain.RegisteredUser) (err error) {
	fmt.Printf("Updating user %s %s\n", user.FirstName, user.Connections)
	filter := bson.M{"_id": user.Id}
	replacementObj := user
	_, err = store.users.ReplaceOne(context.TODO(), filter, replacementObj)
	fmt.Printf("Updated \n")
	if err != nil {
		return err
	}
	return nil
}

func (store *UserMongoDBStore) filterOne(filter interface{}) (user *domain.RegisteredUser, err error) {
	result := store.users.FindOne(context.TODO(), filter)
	err = result.Decode(&user)
	return
}

func (store *UserMongoDBStore) GetBasicInfo() ([]*domain.RegisteredUser, error) {
	projection := bson.D{{"first_name", 1}, {"last_name", 1}}
	opts := options.Find().SetProjection(projection)
	cursor, err := store.users.Find(context.TODO(), bson.D{}, opts)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

//func (store *UserMongoDBStore) FindByFilter(nameFilter string) ([]*domain.RegisteredUser, error) {
//	filter := bson.D{
//		{"first_name", primitive.Regex{Pattern: nameFilter, Options: "i"}},
//		{"$or", []interface{}{
//			bson.D{{"last_name", primitive.Regex{Pattern: nameFilter, Options: "i"}}},
//		}},
//	}
//	return store.filter(filter)
//}

func (store *UserMongoDBStore) UpdatePersonalInfo(user *domain.RegisteredUser) (primitive.ObjectID, error) {
	result, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"first_name", user.FirstName},
				{"last_name", user.LastName},
				{"gender", user.Gender},
				{"phone_number", user.PhoneNumber},
				{"date_of_birth", user.DateOfBirth},
				{"biography", user.Biography},
				{"email", user.Email},
			}},
		},
	)
	upsertedId := fmt.Sprint(result.UpsertedID)
	objectId, _ := primitive.ObjectIDFromHex(upsertedId)
	return objectId, err
}

func (store *UserMongoDBStore) AddExperience(experience *domain.Experience, userId primitive.ObjectID) error {
	user, _ := store.Get(userId)
	experiences := append(user.Experiences, *experience)
	_, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"experiences", experiences}}},
		},
	)

	return err
}

func decode(cursor *mongo.Cursor) (users []*domain.RegisteredUser, err error) {
	for cursor.Next(context.TODO()) {
		var user domain.RegisteredUser
		err = cursor.Decode(&user)
		if err != nil {
			return
		}
		users = append(users, &user)
	}
	err = cursor.Err()
	return
}
