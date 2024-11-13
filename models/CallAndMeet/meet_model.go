package models

type Meet struct {
    ID            int     `json:"id" gorm:"primaryKey;autoIncrement"`
    ClientName    string  `json:"client_name" gorm:"type:varchar(255)"`
    Address       string  `json:"address" gorm:"type:varchar(255)"`
    ProspectStatus string `json:"prospect_status" gorm:"type:varchar(255)"`
    Latitude      float64 `json:"latitude" gorm:"type:decimal(10,8)"`
    Longitude     float64 `json:"longitude" gorm:"type:decimal(11,8)"`
    MeetResult    string  `json:"meet_result" gorm:"type:text"`
}
