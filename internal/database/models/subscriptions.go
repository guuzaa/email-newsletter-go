package models

import "time"

type Subscription struct {
	ID           string    `gorm:"column:id;not null;primaryKey;type:uuid"`
	Email        string    `gorm:"column:email;not null;unique" form:"email"`
	Name         string    `gorm:"column:name;not null" form:"name"`
	SubscribedAt time.Time `gorm:"column:subscribed_at;not null"`
	Status       string    `gorm:"column:status"`
}

const (
	SubscriptionStatusConfirmed = "confirmed"
	SubscriptionStatusPending   = "pending_confirmation"
)
