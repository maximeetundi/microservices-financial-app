package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func parseObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
