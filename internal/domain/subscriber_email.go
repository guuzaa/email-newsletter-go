package domain

import "github.com/go-playground/validator/v10"

type SubscriberEmail string

func (e SubscriberEmail) String() string {
	return string(e)
}

func SubscriberEmailFrom(email string) (SubscriberEmail, error) {
	if err := validator.New().Var(email, "required,email"); err != nil {
		return "", err
	}
	return SubscriberEmail(email), nil
}
