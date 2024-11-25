package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

const (
    RoleBFA   int = 1
    RoleStaff int = 2
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`         
    Username  string             `bson:"username" json:"username"`
    Email     string             `bson:"email" json:"email"`
    PhoneNum  string             `bson:"phone_num" json:"phone_num"`
    Password  string             `bson:"password" json:"password"`
    Image     string             `bson:"image" json:"image"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
    LastLogin time.Time          `bson:"last_login,omitempty" json:"last_login"`
    Role      int                `bson:"role" json:"role"` 
}

func (u *User) IsBFA() bool {
    return u.Role == RoleBFA
}

func (u *User) IsStaff() bool {
    return u.Role == RoleStaff
}

func (u *User) BeforeCreate() error {
    if u.CreatedAt.IsZero() {
        u.CreatedAt = time.Now()
    }
    return nil
}

func (u *User) TableName() string {
    return "users"
}