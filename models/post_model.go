package models

import "time"

type Post struct{
	ID 			int 		`json:"id" gorm:"primaryKey;autoIncrement"`
	NameBfa 	string		`json:"name_bfa" gorm:"varchar(255)"`
	Text		string		`json:"text" gorm:"varchar(255)"`
	Date		time.Time	`json:"date"`
	CommandAndLike string	`json:"command_and_like" gorm:"varchar(255)"`		 
}