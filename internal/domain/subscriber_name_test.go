package domain_test

import (
	"strings"
	"testing"

	"github.com/guuzaa/email-newsletter/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name    string
		isError bool
	}{
		{name: strings.Repeat("ё", 256), isError: false},
		{name: strings.Repeat("a", 257), isError: true},
		{name: " ", isError: true},
		{name: "", isError: true},
		{name: "name", isError: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name, err := domain.SubscriberNameFrom(tc.name)
			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.name, name.String())
			}
		})
	}
}
