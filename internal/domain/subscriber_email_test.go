package domain_test

import (
	"testing"

	"github.com/guuzaa/email-newsletter/internal/domain"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestSubscriberEmailFrom(t *testing.T) {
	faker := faker.New()
	randomEmail := faker.Internet().Email()
	tests := []struct {
		name          string
		email         string
		expected      string
		expectedError bool
	}{
		{
			name:          "valid email",
			email:         randomEmail,
			expected:      randomEmail,
			expectedError: false,
		},
		{
			name:          "invalid email",
			email:         "invalid-email.com",
			expected:      "",
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := domain.SubscriberEmailFrom(test.email)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result.String())
			}
		})
	}
}
