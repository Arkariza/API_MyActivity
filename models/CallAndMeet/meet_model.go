package models

import (
    "errors"
    "strings"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Meet struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    PhoneNum       string             `bson:"phone_num" json:"phone_num"`
    ClientName     string             `bson:"client_name" json:"client_name"`
    Address        string             `bson:"address" json:"address"`
    ProspectStatus string             `bson:"prospect_status" json:"prospect_status"`
    Latitude       float64            `bson:"latitude" json:"latitude"`
    Longitude      float64            `bson:"longitude" json:"longitude"`
    Date           time.Time          `bson:"date" json:"date"`
    MeetResult     string             `bson:"meet_result" json:"meet_result"`
    CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
    Note           string             `bson:"note" json:"note"`
}

func (m *Meet) Validate() error {
    m.ClientName = strings.TrimSpace(m.ClientName)
    if m.ClientName == "" {
        return errors.New("client name is required and cannot be empty")
    }
    if len(m.ClientName) < 2 || len(m.ClientName) > 100 {
        return errors.New("client name must be between 2 and 100 characters")
    }

    m.Address = strings.TrimSpace(m.Address)
    if m.Address == "" {
        return errors.New("address is required and cannot be empty")
    }

    m.PhoneNum = strings.TrimSpace(m.PhoneNum)
    if m.PhoneNum == "" {
        return errors.New("phone number is required")
    }

    if m.Latitude < -90 || m.Latitude > 90 {
        return errors.New("invalid latitude, must be between -90 and 90")
    }
    if m.Longitude < -180 || m.Longitude > 180 {
        return errors.New("invalid longitude, must be between -180 and 180")
    }

    validStatuses := map[string]bool{
        "potential": true,
        "active":    true,
        "inactive":  true,
    }
    if m.ProspectStatus != "" && !validStatuses[m.ProspectStatus] {
        return errors.New("invalid prospect status, must be one of: potential, active, inactive")
    }

    return nil
}

func (m *Meet) BeforeCreate() {
    m.CreatedAt = time.Now()
    if m.ProspectStatus == "" {
        m.ProspectStatus = "potential"
    }
    if m.Note == "" {
        m.Note = "No additional notes provided."
    }
}

func (u *Meet) TableName() string {
    return "meet"
}