package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Call struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ClientName     string             `bson:"client_name" json:"client_name"`
    PhoneNum       string             `bson:"phonenum" json:"phone_num"`
    ProspectStatus string             `bson:"prospect_status" json:"prospect_status"`
    Date           time.Time          `bson:"date" json:"date"`
    Note           string             `bson:"note" json:"note"`
    CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
    CallResult     string             `bson:"call_result" json:"call_result"`
}

func (c *Call) Validate() error {
    c.ClientName = strings.TrimSpace(c.ClientName)
    if c.ClientName == "" {
        return errors.New("client name is required and cannot be empty")
    }
    if len(c.ClientName) < 2 || len(c.ClientName) > 100 {
        return errors.New("client name must be between 2 and 100 characters")
    }

    validStatuses := map[string]bool{
        "new":          true,
        "in_progress":  true,
        "contacted":    true,
        "qualified":    true,
        "unqualified":  true,
        "follow_up":    true,
    }
    if c.ProspectStatus != "" && !validStatuses[c.ProspectStatus] {
        return errors.New("invalid prospect status, must be one of: new, in_progress, contacted, qualified, unqualified, follow_up")
    }

    return nil
}

func (c *Call) BeforeCreate() {
    c.Date = time.Now()
    if c.ProspectStatus == "" {
        c.ProspectStatus = "new"
    }
    if c.Note == "" {
        c.Note = "No additional notes provided."
    }
    if c.CallResult == "" {
        c.CallResult = "Pending"
    }
}