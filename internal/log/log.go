package log

import "time"

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Entry struct {
	Time    time.Time
	Level   Level
	Message string
	Owner   any
}
