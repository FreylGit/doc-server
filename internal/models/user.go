package models

import "time"

type User struct {
	Id        int64
	Login     string
	Password  string
	CreatedAt time.Time
}
