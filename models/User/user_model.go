package models

import "time"

const (
    RoleBFA   int8 = 1
    RoleStaff int8 = 2
)

type User struct {
    ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
    Username  string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"`
    Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
    PhoneNum  string    `gorm:"type:varchar(15);index" json:"phone_num"`
    Password  string    `gorm:"type:varchar(255);not null" json:"password"`
    Image     string    `gorm:"type:varchar(255)" json:"image"`
    CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
    LastLogin time.Time `gorm:"type:timestamp" json:"last_login"`
    Role      int8      `gorm:"type:smallint;default:2;not null" json:"role"`
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