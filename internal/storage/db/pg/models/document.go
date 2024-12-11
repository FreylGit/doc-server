package models

import "time"

type Document struct {
	Id        int64
	Name      string
	Mime      string
	IsPublic  bool
	IsFile    bool
	UserId    int64
	CreatedAt time.Time
	Grant     []Grant
}
