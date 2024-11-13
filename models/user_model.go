package models

type User struct{
	Id			int64	`gorm:"primaryKey" json:"id"`
	UserName 	string	`gorm:"type:varchar(255)" json:"UserName"`
	Email 		string	`gorm:"type:varchar(255)" json:"Email"`
	NumPhone 	int16	`gorm:"type:integer(15)" json:"NumPhone"`
	Password 	string	`gorm:"type:varchar(8)" json:"Password"`
	Img 		string	`gorm:"type:varchar(255)" json:"Image"`
	CreateAt	string	`gorm:"type:date" json:"CreateAt"`
	LogAt      	string	`gorm:"type:date" json:"LogAt"`
	Role      	int16	`gorm:"type:integer(1)" json:"Role"`
}

func (u *User) IsBFA() bool {
	return u.Role == 1
}
func (u *User) IsStaf() bool {
	return u.Role == 2
}