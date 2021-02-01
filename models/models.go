package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Article is a data structure that represents a single record of the articles collection in MongoDB
type Article struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Likes int64              `json:"likes" bson:"likes"`
}
