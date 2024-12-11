package models

import "time"

type Grant struct {
	Id         int64
	Login      string
	DocumentId int64
	Permission string
	CreatedAt  time.Time
}
