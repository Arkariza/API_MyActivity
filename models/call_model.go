package models

import "time"

type Call struct {
    ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
    ClientName     string    `json:"client_name" gorm:"type:varchar(255)"`
    Numphone       int       `json:"numphone"`
    ProspectStatus string    `json:"prospect_status" gorm:"type:varchar(255)"`
    Date           time.Time `json:"date"`
    Note           string    `json:"note" gorm:"type:text"`
    CallResult     string    `json:"call_result" gorm:"type:text"`
}
