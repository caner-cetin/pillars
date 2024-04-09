package db

import (
	"pillars-backend/src/constants"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetCollection(collName string) *mongo.Collection {
	return constants.MONGODB.Database(constants.DATABASE_NAME).Collection(collName)
}
