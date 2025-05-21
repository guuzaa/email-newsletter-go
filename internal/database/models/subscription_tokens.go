package models

type SubscriptionTokens struct {
	SubscriptionToken string `gorm:"column:subscription_token;primaryKey;not null"`
	SubscriptionID    string `gorm:"column:subscription_id;not null;type:uuid"`
}
