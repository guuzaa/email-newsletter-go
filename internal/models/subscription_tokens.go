package models

type SubscriptionTokens struct {
	SubscriptionToken string `gorm:"subscription_token;primaryKey;not null"`
	SubscriptionID    string `gorm:"subscription_id;not null;type:uuid"`
}
