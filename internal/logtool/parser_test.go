package logtool

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestCountLogLevels(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input  io.Reader
		result map[string]int
	}{
		"happy case": {
			input: strings.NewReader(`
				2026-06-25 09:00:00 INFO Starting payment-service
				2026-06-25 09:00:01 INFO Loading configuration
				2026-06-25 09:00:02 INFO Configuration loaded
				2026-06-25 09:00:03 INFO Connecting to PostgreSQL
				2026-06-25 09:00:09 WARN Cache miss for customer:1001
				2026-06-25 09:00:14 WARN Slow database query (340ms)
				2026-06-25 09:00:15 INFO GET /payments/123 200
				2026-06-25 09:00:16 INFO GET /payments/124 200
				2026-06-25 09:00:17 INFO GET /payments/125 200
				2026-06-25 09:00:18 ERROR Failed to process payment id=125`),
			result: map[string]int{
				"INFO":  7,
				"WARN":  2,
				"ERROR": 1,
			},
		},
		"empty string": {
			input: strings.NewReader(""),
			result: map[string]int{
				"INFO":  0,
				"WARN":  0,
				"ERROR": 0,
			},
		},
		"invalid line": {
			input: strings.NewReader(`
				INVALID
				2026-06-25 09:00:00 INFO Starting payment-service
				INVALID
				2026-06-25 09:00:01 INFO Loading configuration
				2026-06-25 09:00:02 INFO Configuration loaded
				2026-06-25 09:00:03 INFO Connecting to PostgreSQL
				INVALID`),
			result: map[string]int{
				"INFO":  4,
				"ERROR": 0,
				"WARN":  0,
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := CountLogLevels(t.Context(), test.input)
			if err != nil {
				t.Fatalf("CountLogLevels() error = %v", err)
			}

			if !reflect.DeepEqual(got, test.result) {
				t.Fatalf("CountLogLevels() = %v, want %v", got, test.result)
			}
		})
	}
}
