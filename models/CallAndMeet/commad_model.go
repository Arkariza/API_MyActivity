package models

import "time"

type Command struct {
    ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
    BFAName         string    `gorm:"type:varchar(255);not null" json:"bfa_name"`
    Description     string    `gorm:"type:text" json:"description"`
    Date            time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"date"`
    CommandAndLink  string    `gorm:"type:text" json:"command_and_link"`
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