package models

import "time"

type Subscription struct {
	ID           string    `gorm:"id;not null;primaryKey"`
	Email        string    `gorm:"email;not null;unique" form:"email"`
	Name         string    `gorm:"name;not null" form:"name"`
	SubscribedAt time.Time `gorm:"subscribed_at;not null"`
}
