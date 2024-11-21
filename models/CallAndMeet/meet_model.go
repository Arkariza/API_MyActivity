package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Meet struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          int                `json:"user_id" bson:"user_id"`
	ClientName      string             `json:"client_name" bson:"client_name" validate:"required,min=2,max=100"`
	Address         string             `json:"address" bson:"address" validate:"required"`
	ProspectStatus  string             `json:"prospect_status" bson:"prospect_status"`
	Latitude        float64            `json:"latitude" bson:"latitude" validate:"required"`
	Longitude       float64            `json:"longitude" bson:"longitude" validate:"required"`
	MeetResult      string             `json:"meet_result" bson:"meet_result"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}

func (m *Meet) Validate() error {
	m.ClientName = strings.TrimSpace(m.ClientName)
	if m.ClientName == "" {
		return errors.New("client name is required")
	}

	m.Address = strings.TrimSpace(m.Address)
	if m.Address == "" {
		return errors.New("address is required")
	}

	if m.Latitude < -90 || m.Latitude > 90 {
		return errors.New("invalid latitude. Must be between -90 and 90")
	}

	if m.Longitude < -180 || m.Longitude > 180 {
		return errors.New("invalid longitude. Must be between -180 and 180")
	}

	validStatuses := map[string]bool{
		"potential": true,
		"active":    true,
		"inactive":  true,
		"":          true, 
	}

	if !validStatuses[m.ProspectStatus] {
		return errors.New("invalid prospect status. Must be one of: potential, active, inactive")
	}

	return nil
}

func (m *Meet) BeforeCreate() {
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
		if m.ProspectStatus == "" {
		m.ProspectStatus = "potential"
	}
}

func (m *Meet) TableName() string {
    return "meet"
}