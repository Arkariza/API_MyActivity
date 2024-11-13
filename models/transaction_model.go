package models

import "time"

type Transaction struct {
    ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
    PhoneNumber  string    `gorm:"type:varchar(15)" json:"phone_number"`
    Priority     string    `gorm:"type:varchar(50)" json:"priority"`
    PolicyNumber int32     `gorm:"type:integer" json:"policy_number"`
    Information  string    `gorm:"type:text" json:"information"`
    CreatedAt    time.Time `gorm:"type:timestamp;not null" json:"created_at"`
    SubmitDate   time.Time `gorm:"type:timestamp" json:"submit_date"`
    ClientName   string    `gorm:"type:varchar(255)" json:"client_name"`
    BFAId        uint      `gorm:"type:integer" json:"bfa_id"`
    ReferralID   uint      `gorm:"type:integer" json:"referral_id"`
    LeadType     string    `gorm:"type:varchar(10)" json:"lead_type"`
    Status       string    `gorm:"type:varchar(50)" json:"status"`
}