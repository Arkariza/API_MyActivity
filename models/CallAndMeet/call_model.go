package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Call struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ClientName     string    `json:"client_name" bson:"type:varchar(255)"`
    Numphone       int       `json:"numphone"`
    ProspectStatus string    `json:"prospect_status" bson:"type:varchar(255)"`
    Date           time.Time `json:"date"`
    Note           string    `json:"note" bson:"type:text"`
    CallResult     string    `json:"call_result" bson:"type:text"`
}
