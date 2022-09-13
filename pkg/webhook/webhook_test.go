package pkg_webhook

import (
	"testing"

	klockv1 "klock/apis/klock/v1"

	v1 "k8s.io/api/authentication/v1"

	"github.com/stretchr/testify/assert"
)

func Test_exclude(t *testing.T) {
	type test struct {
		name      string
		exclusive klockv1.Exclusive
		user      v1.UserInfo

		expected bool
	}

	tests := []test{
		{
			name:      "not exclusive",
			exclusive: klockv1.Exclusive{},
			user: v1.UserInfo{
				Username: "robert",
				UID:      "rrr-55-rr-rr",
			},

			expected: false,
		},
		{
			name: "exclusive, just name",
			exclusive: klockv1.Exclusive{
				Name: "robert",
			},
			user: v1.UserInfo{
				Username: "robert",
				UID:      "rrr-55-rr-rr",
			},

			expected: true,
		},
		{
			name: "exclusive, just uid",
			exclusive: klockv1.Exclusive{
				UID: "rrr-55-rr-rr",
			},
			user: v1.UserInfo{
				Username: "robert",
				UID:      "rrr-55-rr-rr",
			},

			expected: true,
		},
		{
			name: "exclusive, both",
			exclusive: klockv1.Exclusive{
				Name: "robert",
				UID:  "rrr-55-rr-rr",
			},
			user: v1.UserInfo{
				Username: "robert",
				UID:      "rrr-55-rr-rr",
			},

			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := exclude(tt.exclusive, tt.user)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
