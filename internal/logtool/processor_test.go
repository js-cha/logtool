package logtool

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestProcessLogFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		content []byte
		result  map[string]int
	}{
		"happy case": {
			input: "api.log",
			content: []byte(`
2026-06-25 09:00:00 INFO Starting payment-service
2026-06-25 09:00:01 WARN Cache miss
2026-06-25 09:00:02 ERROR Failed to process payment`),
			result: map[string]int{
				"INFO":  1,
				"WARN":  1,
				"ERROR": 1,
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			path := filepath.Join(dir, test.input)

			if err := os.WriteFile(path, test.content, 0o644); err != nil {
				t.Fatalf("write temp log file: %v", err)
			}

			got, err := ProcessLogFile(t.Context(), path)
			if err != nil {
				t.Fatalf("CountLogLevels() error = %v", err)
			}

			if !reflect.DeepEqual(got, test.result) {
				t.Fatalf("CountLogLevels() = %v, want %v", got, test.result)
			}
		})
	}
}

func TestProcessDirConcurrent(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   []string
		content [][]byte
		result  map[string]int
	}{
		"happy case": {
			input: []string{"api.log", "payments.log", "auth.log"},
			content: [][]byte{
				[]byte(`
2026-06-25 09:00:00 INFO Starting payment-service
2026-06-25 09:00:01 WARN Cache miss
2026-06-25 09:00:02 ERROR Failed to process payment`),
				[]byte(`
2026-06-25 09:00:00 INFO Starting payment-service
2026-06-25 09:00:01 WARN Cache miss
2026-06-25 09:00:02 ERROR Failed to process payment`),
				[]byte(`
2026-06-25 09:00:00 INFO Starting payment-service
2026-06-25 09:00:01 WARN Cache miss
2026-06-25 09:00:02 ERROR Failed to process payment`),
			},
			result: map[string]int{
				"INFO":  3,
				"WARN":  3,
				"ERROR": 3,
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()

			for i := 0; i < len(test.input); i++ {
				path := filepath.Join(dir, test.input[i])

				if err := os.WriteFile(path, test.content[i], 0o644); err != nil {
					t.Fatalf("write temp log file: %v", err)
				}
			}

			got, err := ProcessDirConcurrent(t.Context(), dir, 4)
			if err != nil {
				t.Fatalf("CountLogLevels() error = %v", err)
			}

			if !reflect.DeepEqual(got, test.result) {
				t.Fatalf("CountLogLevels() = %v, want %v", got, test.result)
			}
		})
	}
}
