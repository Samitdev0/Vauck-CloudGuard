package middleware

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuditOutcome(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		expected   string
	}{
		{
			name:       "successful request",
			statusCode: 200,
			err:        nil,
			expected:   "success",
		},
		{
			name:       "created request",
			statusCode: 201,
			err:        nil,
			expected:   "success",
		},
		{
			name:       "unauthorized request",
			statusCode: 401,
			err:        nil,
			expected:   "denied",
		},
		{
			name:       "forbidden request",
			statusCode: 403,
			err:        nil,
			expected:   "denied",
		},
		{
			name:       "bad request",
			statusCode: 400,
			err:        nil,
			expected:   "failure",
		},
		{
			name:       "not found",
			statusCode: 404,
			err:        nil,
			expected:   "failure",
		},
		{
			name:       "server error",
			statusCode: 500,
			err:        nil,
			expected:   "failure",
		},
		{
			name:       "handler returned error",
			statusCode: 200,
			err:        echo.NewHTTPError(500, "internal error"),
			expected:   "failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := auditOutcome(tt.statusCode, tt.err)

			if got != tt.expected {
				t.Fatalf(
					"expected outcome %q, got %q",
					tt.expected,
					got,
				)
			}
		})
	}
}
