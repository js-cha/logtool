package logtool

type Result struct {
	Counts map[string]int
	Err    error
}

type LogEntry struct {
	Level string
}

type LogLevel string

const (
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
)
