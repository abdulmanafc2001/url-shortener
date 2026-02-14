package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
	}
}

// TODO : add improvments in the logging

func (l *Logger) formatMessage(msg string, fields map[string]interface{}) string {
	// Get caller information
	_, file, line, _ := runtime.Caller(2)
	// Get just the file name without the full path
	parts := strings.Split(file, "/")
	fileName := parts[len(parts)-1]

	// Format fields
	fieldStr := ""
	if len(fields) > 0 {
		pairs := make([]string, 0, len(fields))
		for k, v := range fields {
			pairs = append(pairs, fmt.Sprintf("%s=%v", k, v))
		}
		fieldStr = " | " + strings.Join(pairs, " | ")
	}

	// Return formatted message with timestamp, file location, and fields
	return fmt.Sprintf("[%s:%d] %s%s", fileName, line, msg, fieldStr)
}

func (l *Logger) Info(msg string, fields map[string]any) {
	l.infoLogger.Println(l.formatMessage(msg, fields))
}

func (l *Logger) Error(msg string, fields map[string]any) {
	l.errorLogger.Println(l.formatMessage(msg, fields))
}

func (l *Logger) Debug(msg string, fields map[string]any) {
	l.debugLogger.Println(l.formatMessage(msg, fields))
}
