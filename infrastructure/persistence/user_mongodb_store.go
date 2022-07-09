package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_UserService/tracer"
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

func (store *UserMongoDBStore) GetActiveById(ctx context.Context, id primitive.ObjectID) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"_id": id}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) GetAllActive(ctx context.Context) ([]*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.D{{}}
	return store.filter(ctx, filter)
}

func (store *UserMongoDBStore) GetActiveByUsername(ctx context.Context, username string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"username": username}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) GetByUsername(ctx context.Context, username string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"username": username}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) GetActiveByEmail(ctx context.Context, email string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"email": email}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) GetByEmail(ctx context.Context, email string) (*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	filter := bson.M{"email": email}
	return store.filterOne(ctx, filter)
}

func (store *UserMongoDBStore) Insert(ctx context.Context, user *domain.RegisteredUser) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *UserMongoDBStore) UpdateIsActive(ctx context.Context, email string) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.D{{"$set", bson.D{{"is_active", true}}}},
	)
	return err
}

func (store *UserMongoDBStore) DeleteAll(ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	store.users.DeleteMany(ctx, bson.D{{}})
}

func (store *UserMongoDBStore) filter(ctx context.Context, filter interface{}) ([]*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	cursor, err := store.users.Find(ctx, filter)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *UserMongoDBStore) Update(ctx context.Context, user *domain.RegisteredUser) (err error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Printf("Updating user %s %s\n", user.FirstName, user.Connections)
	filter := bson.M{"_id": user.Id}
	replacementObj := user
	_, err = store.users.ReplaceOne(ctx, filter, replacementObj)
	fmt.Printf("Updated \n")
	if err != nil {
		return err
	}
	return nil
}

func (store *UserMongoDBStore) filterOne(ctx context.Context, filter interface{}) (user *domain.RegisteredUser, err error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result := store.users.FindOne(ctx, filter)
	err = result.Decode(&user)
	return
}

func (store *UserMongoDBStore) GetBasicInfo(ctx context.Context) ([]*domain.RegisteredUser, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	projection := bson.D{{"first_name", 1}, {"last_name", 1}, {"is_active", true}}
	opts := options.Find().SetProjection(projection)
	cursor, err := store.users.Find(ctx, bson.D{}, opts)
	defer cursor.Close(ctx)

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *UserMongoDBStore) UpdatePersonalInfo(ctx context.Context, user *domain.RegisteredUser) (primitive.ObjectID, error) {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	result, err := store.users.UpdateOne(
		ctx,
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

func (store *UserMongoDBStore) AddExperience(ctx context.Context, experience *domain.Experience, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, _ := store.GetActiveById(ctx, userId)
	experiences := append(user.Experiences, *experience)
	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"experiences", experiences}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) AddEducation(ctx context.Context, education *domain.Education, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, _ := store.GetActiveById(ctx, userId)
	educations := append(user.Educations, *education)
	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"educations", educations}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) ChangeAccountPrivacy(ctx context.Context, isPrivate bool, username string) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, _ := store.GetByUsername(ctx, username)
	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"is_private", isPrivate}}},
		},
	)
	return err
}

func (store *UserMongoDBStore) AddSkill(ctx context.Context, skill string, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, _ := store.GetActiveById(ctx, userId)
	skillExists := util.ContainsStr(user.Skills, skill)
	if skillExists {
		err := errors.New("skill already exists")
		return err
	}
	skills := append(user.Skills, skill)
	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"skills", skills}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) RemoveSkill(ctx context.Context, removeSkill string, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := store.GetActiveById(ctx, userId)
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
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"skills", skills}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) AddInterest(ctx context.Context, companyId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, _ := store.GetActiveById(ctx, userId)
	interestExists := util.ContainsId(user.Interests, companyId)
	if interestExists {
		err := errors.New("interest already added")
		return err
	}
	interests := append(user.Interests, companyId)
	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"interests", interests}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) DeleteExperience(ctx context.Context, experienceId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, _ := store.GetActiveById(ctx, userId)
	experiences, errDel := util.DeleteExperience(user.Experiences, experienceId)
	if errDel != nil {
		return errDel
	}
	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"experiences", experiences}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) DeleteEducation(ctx context.Context, educationId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, _ := store.GetActiveById(ctx, userId)
	educations, errDel := util.DeleteEducation(user.Educations, educationId)
	if errDel != nil {
		return errDel
	}
	_, err := store.users.UpdateOne(
		ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"educations", educations}}},
		},
	)

	return err
}

func (store *UserMongoDBStore) RemoveInterest(ctx context.Context, companyId primitive.ObjectID, userId primitive.ObjectID) error {
	span := tracer.StartSpanFromContext(ctx, "DB Get")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	user, err := store.GetActiveById(ctx, userId)
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
		ctx,
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
