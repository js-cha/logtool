package logtool

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

func parsLogEntry(line string) (LogEntry, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return LogEntry{}, fmt.Errorf("invalid log line: %q", line)
	}
	level := fields[2]
	switch level {
	case "INFO", "WARN", "ERROR":
		return LogEntry{Level: level}, nil
	default:
		return LogEntry{}, fmt.Errorf("unknown log level: %q", level)
	}
}

func CountLogLevels(ctx context.Context, reader io.Reader) (map[string]int, error) {
	logLevelMap := map[LogLevel]int{
		"INFO":  0,
		"WARN":  0,
		"ERROR": 0,
	}

	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024) // 64KB intial size
	scanner.Buffer(buf, 1024*1024)  // 1MB cap

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		text := scanner.Text()

		logLevel, _ := parsLogEntry(text)
		if _, ok := logLevelMap[LogLevel(logLevel.Level)]; ok {
			logLevelMap[LogLevel(logLevel.Level)]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	out := map[string]int{
		"INFO":  int(logLevelMap[InfoLevel]),
		"WARN":  int(logLevelMap[WarnLevel]),
		"ERROR": int(logLevelMap[ErrorLevel]),
	}

	return out, nil
}

func PrintLogLevels(m map[string]int) {
	levels := []string{"INFO", "WARN", "ERROR"}
	for _, level := range levels {
		fmt.Printf("%s: %d\n", level, m[level])
	}
}
