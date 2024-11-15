package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct {
	primitive.ObjectID 	  `bson:"_id,omitempty" json:"id"`
	NumPhone    int16     `bson:"integer(15)" json:"NumPhone"`
	Priority    string    `bson:"varchar(10)" json:"Priority"`
	Latitude    float64   `bson:"type:decimal(10,6)" json:"latitude"`
	Longitude   float64   `bson:"type:decimal(10,6)" json:"longitude"`
	CreateAt    time.Time `bson:"type:timestamp;not null" json:"created_at"`
	DateSubmit  time.Time `bson:"type:timestamp" json:"date_submit"`
	ClientName  string    `bson:"varchar(50)" json:"ClientName"`
	IdBFA       uint      `bson:"type:integer" json:"bfa_id"`
	IdRefeal    uint      `bson:"type:integer" json:"referral_id"`
	TypeLead    string    `bson:"varchar(10)" json:"TypeLead"`
	NoPolicy    int32     `bson:"integer" json:"NoPolicy"`
	Information string    `bson:"text" json:"Information"`
}