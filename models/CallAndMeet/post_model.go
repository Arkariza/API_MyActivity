package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct{
	ID 			primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NameBfa 	string		`json:"name_bfa" bson:"varchar(255)"`
	Text		string		`json:"text" bson:"varchar(255)"`
	Date		time.Time	`json:"date"`
	CommandAndLike string	`json:"command_and_like" bson:"varchar(255)"`		 
}