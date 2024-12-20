package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct {
    ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
    UserID      primitive.ObjectID  `bson:"user_id" json:"user_id"`
    NumPhone    string              `bson:"numphone" json:"numPhone" binding:"required"`
    Priority    string              `bson:"priority" json:"priority" binding:"required"`
    Latitude    float64             `bson:"latitude" json:"latitude" binding:"required"`
    Longitude   float64             `bson:"longitude" json:"longitude" binding:"required"`
    CreateAt    time.Time           `bson:"created_at" json:"createdAt"`
    DateSubmit  time.Time           `bson:"date_submit,omitempty" json:"dateSubmit"`
    ClientName  string              `bson:"clientname" json:"clientName" binding:"required"`
    TypeLead    string              `bson:"type_lead" json:"typeLead"`
    NoPolicy    int32               `bson:"no_policy,omitempty" json:"noPolicy"`
    Information string              `bson:"information" json:"information"`
    Status      string              `bson:"status" json:"status" binding:"required"`
}

const (
    StatusPending = "Pending"
    StatusWin     = "Win"
    StatusLose    = "Lose"
    StatusOpen    = "Open"

    TypeReferral = "Reff"
    TypeSelf     = "Self"
)

func (l *Lead) ValidateStatus() bool {
    return l.Status == StatusPending || l.Status == StatusWin || l.Status == StatusLose || l.Status == StatusOpen
}

func (l *Lead) ValidateTypeLead() bool {
    return l.TypeLead == TypeReferral || l.TypeLead == TypeSelf
}

func (l *Lead) BeforeCreate() {
    if l.ID.IsZero() {
        l.ID = primitive.NewObjectID()
    }
    if l.CreateAt.IsZero() {
        l.CreateAt = time.Now()
    }
}

func (l *Lead) TableName() string {
    return "leads"
}

type LeadInput struct {
    NumPhone    string  `json:"numPhone" binding:"required"`
    Priority    string  `json:"priority" binding:"required"`
    Latitude    float64 `json:"latitude" binding:"required"`
    Longitude   float64 `json:"longitude" binding:"required"`
    ClientName  string  `json:"clientName" binding:"required"`
    TypeLead    string  `json:"typeLead" binding:"required"`
    NoPolicy    int32   `json:"noPolicy"`
    Information string  `json:"information"`
    Status      string  `json:"status" binding:"required"`
}
