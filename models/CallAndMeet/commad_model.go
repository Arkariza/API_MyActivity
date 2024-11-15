package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Command struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    BFAName         string    `bson:"type:varchar(255);not null" json:"bfa_name"`
    Description     string    `bson:"type:text" json:"description"`
    Date            time.Time `bson:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"date"`
    CommandAndLink  string    `bson:"type:text" json:"command_and_link"`
}

func (c *Command) TableName() string {
    return "commands"
}

func (c *Command) BeforeCreate() error {
    if c.Date.IsZero() {
        c.Date = time.Now()
    }
    return nil
}