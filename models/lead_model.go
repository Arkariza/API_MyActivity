package models

import "time"

type Lead struct {
	Id          int       `gorm:"primarykey:autoincrement" json:"Id"`
	NumPhone    int16     `gorm:"integer(15)" json:"NumPhone"`
	Priority    string    `gorm:"varchar(10)" json:"Priority"`
	Latitude    float64   `gorm:"type:decimal(10,6)" json:"latitude"`
	Longitude   float64   `gorm:"type:decimal(10,6)" json:"longitude"`
	CreateAt    time.Time `gorm:"type:timestamp;not null" json:"created_at"`
	DateSubmit  time.Time `gorm:"type:timestamp" json:"date_submit"`
	ClientName  string    `gorm:"varchar(50)" json:"ClientName"`
	IdBFA       uint      `gorm:"type:integer" json:"bfa_id"`
	IdRefeal    uint      `gorm:"type:integer" json:"referral_id"`
	TypeLead    string    `gorm:"varchar(10)" json:"TypeLead"`
	NoPolicy    int32     `gorm:"integer" json:"NoPolicy"`
	Information string    `gorm:"text" json:"Information"`
}