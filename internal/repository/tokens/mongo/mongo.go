package mongo

import (
	"context"

	"github.com/ikarizxc/http-server/pkg/tokens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenStorage struct {
	client *mongo.Client
}

const (
	dbName   = "authentication-service"
	collName = "refresh-tokens"
)

func NewTokensStorage(client *mongo.Client) *TokenStorage {
	return &TokenStorage{client}
}

func (ts *TokenStorage) Disconnect() error {
	return ts.client.Disconnect(context.TODO())
}

func (ts *TokenStorage) WriteRefreshToken(id int, refreshToken string) error {
	coll := ts.client.Database(dbName).Collection(collName)

	doc := bson.M{
		"id":           id,
		"refreshToken": refreshToken,
	}

	_, err := coll.InsertOne(context.TODO(), doc)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStorage) UpdateRefreshToken(id int, refreshToken string) error {
	coll := ts.client.Database(dbName).Collection(collName)

	filter := bson.M{"id": id}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "refreshToken", Value: refreshToken}}}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TokenStorage) ReadRefreshToken(id int) (string, error) {
	coll := ts.client.Database(dbName).Collection(collName)

	var readToken tokens.Token

	filter := bson.M{"id": id}

	err := coll.FindOne(context.TODO(), filter).Decode(&readToken)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", err
		}

		return "", err
	}

	return readToken.RefreshToken, nil
}
