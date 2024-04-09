package db

import (
	"context"
	"pillars-backend/src/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes() {
	// no need to worry about running this a dozens of times, it will not panic, and we will not handle the error
	// if index already exists, it will not be created again.
	GetCollection(constants.EARTHQUAKE_COLLECTION_NAME).
		Indexes().
		CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		},
		{
			Keys:    bson.D{{Key: "$**", Value: "text"}},
			Options: options.Index().SetName("wildcardTextIndex"),
		},
		{
			Keys: bson.D{{Key: "time", Value: 1}},
		},
	})
}
