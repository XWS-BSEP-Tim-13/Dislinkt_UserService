package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	util "github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/util"
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

func (store *UserMongoDBStore) GetActiveById(id primitive.ObjectID) (*domain.RegisteredUser, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetAllActive() ([]*domain.RegisteredUser, error) {
	filter := bson.D{{"is_active", true}}
	return store.filter(filter)
}

func (store *UserMongoDBStore) GetActiveByUsername(username string) (*domain.RegisteredUser, error) {
	filter := bson.M{"username": username, "is_active": true}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetByUsername(username string) (*domain.RegisteredUser, error) {
	filter := bson.M{"username": username}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetActiveByEmail(email string) (*domain.RegisteredUser, error) {
	filter := bson.M{"email": email, "is_active": true}
	return store.filterOne(filter)
}

func (store *UserMongoDBStore) GetByEmail(email string) (*domain.RegisteredUser, error) {
	filter := bson.M{"email": email}
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

func (store *UserMongoDBStore) UpdateIsActive(email string) error {
	_, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"email": email},
		bson.D{{"$set", bson.D{{"is_active", true}}}},
	)
	return err
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
	projection := bson.D{{"first_name", 1}, {"last_name", 1}, {"is_active", true}}
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
		bson.M{"_id": user.Id, "is_active": true},
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
	user, _ := store.GetActiveById(userId)
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

func (store *UserMongoDBStore) AddEducation(education *domain.Education, userId primitive.ObjectID) error {
	user, _ := store.GetActiveById(userId)
	educations := append(user.Educations, *education)
	_, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"educations", educations}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) AddSkill(skill string, userId primitive.ObjectID) error {
	user, _ := store.GetActiveById(userId)
	skillExists := util.ContainsStr(user.Skills, skill)
	if skillExists {
		err := errors.New("skill already exists")
		return err
	}
	skills := append(user.Skills, skill)
	_, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"skills", skills}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) RemoveSkill(removeSkill string, userId primitive.ObjectID) error {
	user, err := store.GetActiveById(userId)
	if err != nil {
		return errors.New("no such user")
	}

	skillExists := util.ContainsStr(user.Skills, removeSkill)
	if !skillExists {
		return errors.New("skill doesn't exists")
	}

	var skills []string
	for idx, skill := range user.Skills {
		if skill == removeSkill {
			skills = util.RemoveElement(user.Skills, idx)
			break
		}
	}

	_, err = store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"skills", skills}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) AddInterest(companyId primitive.ObjectID, userId primitive.ObjectID) error {
	user, _ := store.GetActiveById(userId)
	interestExists := util.ContainsId(user.Interests, companyId)
	if interestExists {
		err := errors.New("interest already added")
		return err
	}
	interests := append(user.Interests, companyId)
	_, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"interests", interests}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) DeleteExperience(experienceId primitive.ObjectID, userId primitive.ObjectID) error {
	user, _ := store.GetActiveById(userId)
	experiences, errDel := util.DeleteExperience(user.Experiences, experienceId)
	if errDel != nil {
		return errDel
	}
	_, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"experiences", experiences}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) DeleteEducation(educationId primitive.ObjectID, userId primitive.ObjectID) error {
	user, _ := store.GetActiveById(userId)
	educations, errDel := util.DeleteEducation(user.Educations, educationId)
	if errDel != nil {
		return errDel
	}
	_, err := store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"educations", educations}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) RemoveInterest(companyId primitive.ObjectID, userId primitive.ObjectID) error {
	user, err := store.GetActiveById(userId)
	if err != nil {
		return errors.New("no such user")
	}

	interestExists := util.ContainsId(user.Interests, companyId)
	if !interestExists {
		return errors.New("interest doesn't exists")
	}

	var interests []primitive.ObjectID
	for idx, interest := range user.Interests {
		if interest == companyId {
			interests = util.RemoveIdElement(user.Interests, idx)
			break
		}
	}

	_, err = store.users.UpdateOne(
		context.TODO(),
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"interests", interests}}},
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
