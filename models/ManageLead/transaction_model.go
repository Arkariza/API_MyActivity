package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	ClientName  string              `bson:"clientname" json:"clientName"`
	Status      string              `bson:"status" json:"status"`
	UserID      primitive.ObjectID  `bson:"user_id" json:"user_id"`
}
