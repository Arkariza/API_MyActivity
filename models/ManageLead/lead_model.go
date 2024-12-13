package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Lead struct {
    ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
    NumPhone    int64               `bson:"num_phone" json:"numPhone" binding:"required"`
    Priority    string              `bson:"priority" json:"priority" binding:"required"`
    Latitude    float64             `bson:"latitude" json:"latitude" binding:"required"`
    Longitude   float64             `bson:"longitude" json:"longitude" binding:"required"`
    CreateAt    time.Time           `bson:"created_at" json:"createdAt"`
    DateSubmit  time.Time           `bson:"date_submit,omitempty" json:"dateSubmit"`
    ClientName  string              `bson:"client_name" json:"clientName" binding:"required"`
    IdBFA       uint                `bson:"bfa_id,omitempty" json:"bfaId"`
    IdRefeal    uint                `bson:"referral_id,omitempty" json:"referralId"`
    TypeLead    string              `bson:"type_lead" json:"typeLead" binding:"required"`
    NoPolicy    int32               `bson:"no_policy,omitempty" json:"noPolicy"`
    Information string              `bson:"information" json:"information"`
    Status      string              `bson:"status" json:"status"`
}

const (
    StatusSelf     = "Self"
    StatusReferral = "Referral"
    StatusUnknown  = "Unknown"
)

func (l *Lead) ValidateStatus() bool {
    return l.Status == StatusSelf || l.Status == StatusReferral || l.Status == StatusUnknown
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
    NumPhone    int64     `json:"numPhone" binding:"required"`
    Priority    string    `json:"priority" binding:"required"`
    Latitude    float64   `json:"latitude" binding:"required"`
    Longitude   float64   `json:"longitude" binding:"required"`
    ClientName  string    `json:"clientName" binding:"required"`
    TypeLead    string    `json:"typeLead" binding:"required"`
    NoPolicy    int32     `json:"noPolicy"`
    Information string    `json:"information"`
}