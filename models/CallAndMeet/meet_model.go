package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Meet struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ClientName      string  `json:"client_name" bson:"type:varchar(255)"`
    Address         string  `json:"address" bson:"type:varchar(255)"`
    ProspectStatus  string  `json:"prospect_status" bson:"type:varchar(255)"`
    Latitude        float64 `json:"latitude" bson:"type:decimal(10,8)"`
    Longitude       float64 `json:"longitude" bson:"type:decimal(11,8)"`
    MeetResult      string  `json:"meet_result" bson:"type:text"`
}
