package books

import "time"

type BookInput struct {
	ID         	string	`gorm:"primaryKey"`
	Name       	string
	Year       	int
	Author     	string
	Summary    	string
	Publisher  	string
	PageCount  	int
	ReadPage   	int
	Finished   	bool
	Reading    	bool
	InsertedAt 	time.Time
	UpdatedAt	time.Time
}