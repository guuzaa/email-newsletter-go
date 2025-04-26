package domain

import (
	"errors"
	"strings"
)

const maxLength = 256

type SubscriberName string

func (n SubscriberName) String() string {
	return string(n)
}

func SubscriberNameFrom(name string) (SubscriberName, error) {
	const forbiddenChars = "(){}<>[]\\/\""

	isEmptyOrWhitespace := len(strings.TrimSpace(name)) == 0
	isTooLong := len([]rune(name)) > maxLength
	containsForbiddenChars := strings.ContainsAny(name, forbiddenChars)
	if isEmptyOrWhitespace || isTooLong || containsForbiddenChars {
		return "", errors.New("invalid name")
	}
	return SubscriberName(name), nil
}
