package main

var defaultLogger *Logger = NewLogger()

const DefaultLevel int = 1

type Logger struct {
	Level int
}

func (l *Logger) SetLevel(level int) {
	l.Level = level
}
func NewLogger() *Logger {
	return &Logger{Level: 0}
}

func init() {
	defaultLogger.SetLevel(DefaultLevel)
}
