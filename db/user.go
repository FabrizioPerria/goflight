package db

import (
	"context"
	"fmt"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStorer interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	DeleteUserById(ctx context.Context, id string) (string, error)
	UpdateUser(ctx context.Context, filter bson.M, values types.UpdateUserParams) (string, error)
	Dropper
}

const (
	userCollection = "users"
)

type MongoDbUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbUserStore(client *mongo.Client) *MongoDbUserStore {
	return &MongoDbUserStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(userCollection),
	}
}

func (db *MongoDbUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format")
	}
	user := &types.User{}
	err = db.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)

	return user, err
}

func (db *MongoDbUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var cursor *mongo.Cursor
	cursor, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	results := make([]*types.User, 0)
	err = cursor.All(ctx, &results)

	return results, err
}

func (db *MongoDbUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	result, err := db.collection.InsertOne(ctx, user)
	user.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return user, err
}

func (db *MongoDbUserStore) DeleteUserById(ctx context.Context, id string) (string, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	res, err := db.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil || res.DeletedCount == 0 {
		return "", fmt.Errorf("user %s not found", id)
	}

	return id, nil
}

func (db *MongoDbUserStore) Drop(ctx context.Context) error {
	err := db.collection.Drop(ctx)
	return err
}

func (db *MongoDbUserStore) UpdateUser(ctx context.Context, filter bson.M, values types.UpdateUserParams) (string, error) {
	result, err := db.collection.UpdateOne(ctx, filter, bson.M{"$set": values})
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("user not found")
	}
	return "", nil
}
