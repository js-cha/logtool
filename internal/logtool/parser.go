package logtool

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

func CountLogLevels(ctx context.Context, reader io.Reader) (map[string]int, error) {
	logLevelMap := map[string]int{
		"INFO":  0,
		"WARN":  0,
		"ERROR": 0,
	}

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		text := scanner.Text()
		fields := strings.Fields(text)

		if len(fields) < 3 {
			continue
		}

		logLevel := fields[2]
		if _, ok := logLevelMap[logLevel]; ok {
			logLevelMap[logLevel]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logLevelMap, nil
}

func PrintLogLevels(m map[string]int) {
	levels := []string{"INFO", "WARN", "ERROR"}
	for _, level := range levels {
		fmt.Printf("%s: %d\n", level, m[level])
	}
}
