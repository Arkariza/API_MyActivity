package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    PhoneNumber  string             `bson:"type:varchar(15)" json:"phone_number"`
    Priority     string             `bson:"type:varchar(50)" json:"priority"`
    PolicyNumber int32              `bson:"type:integer" json:"policy_number"`
    Information  string             `bson:"type:text" json:"information"`
    CreatedAt    time.Time          `bson:"type:timestamp;not null" json:"created_at"`
    SubmitDate   time.Time          `bson:"type:timestamp" json:"submit_date"`
    ClientName   string             `bson:"type:varchar(255)" json:"client_name"`
    BFAId        uint               `bson:"type:integer" json:"bfa_id"`
    ReferralID   uint               `bson:"type:integer" json:"referral_id"`
    LeadType     string             `bson:"type:varchar(10)" json:"lead_type"`
    Status       string             `bson:"type:varchar(50)" json:"status"`
}