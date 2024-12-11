package models

import (
    "errors"
    "strings"
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title           string             `bson:"title" json:"title"`
    Description     string             `bson:"description" json:"description"`
    Date            time.Time          `bson:"date" json:"date"`
    PostedBy        string             `bson:"posted_by" json:"posted_by"`
    UserRole        int                `bson:"user_role" json:"user_role"`
}

func (c *Comment) Validate() error {
    c.Title = strings.TrimSpace(c.Title)
    c.Description = strings.TrimSpace(c.Description)
    c.PostedBy = strings.TrimSpace(c.PostedBy)

    if c.Title == "" {
        return errors.New("title is required and cannot be empty")
    }
    if len(c.Title) < 2 || len(c.Title) > 255 {
        return errors.New("title must be between 2 and 255 characters")
    }

    if c.PostedBy == "" {
        return errors.New("posted by name is required and cannot be empty")
    }
    if len(c.PostedBy) < 2 || len(c.PostedBy) > 255 {
        return errors.New("posted by name must be between 2 and 255 characters")
    }

    return nil
}

func (c *Comment) BeforeCreate() error {
    if err := c.Validate(); err != nil {
        return err
    }

    if c.Date.IsZero() {
        c.Date = time.Now()
    }

    return nil
}

func CreateComment(title, description, postedBy string, userRole int) (*Comment, error) {
    comment := &Comment{
        Title:          title,
        Description:    description,
        PostedBy:       postedBy,
        UserRole:       userRole,
        Date:           time.Now(),
    }

    err := comment.Validate()
    if err != nil {
        return nil, err
    }

    return comment, nil
}