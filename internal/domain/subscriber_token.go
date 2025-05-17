package domain

import (
	"math/rand"
	"regexp"
	"time"
)

func NewSubscriptionToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 25)
	for i := range token {
		rand.NewSource(time.Now().UnixNano())
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}

func ValidSubscriberToken(token string) bool {
	const tokenPattern = `[a-zA-Z0-9]{25}`
	re := regexp.MustCompile(tokenPattern)
	return re.MatchString(token)
}
