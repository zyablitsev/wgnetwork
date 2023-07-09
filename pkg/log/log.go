package log

import (
	"fmt"
	"io"
	"log"
)

// Logger struct
type Logger struct {
	level int

	errorOut   *log.Logger
	warningOut *log.Logger
	infoOut    *log.Logger
	debugOut   *log.Logger
}

// New contructor
func New(loglevel string, stdout, stderr io.Writer) (*Logger, error) {
	logger := &Logger{
		level:      label2level(loglevel),
		errorOut:   log.New(stderr, "Err: ", log.Ldate|log.Ltime),
		warningOut: log.New(stdout, "Wrn: ", log.Ldate|log.Ltime),
		infoOut:    log.New(stdout, "Inf: ", log.Ldate|log.Ltime),
		debugOut:   log.New(stdout, "Dbg: ", log.Ldate|log.Ltime),
	}

	return logger, nil
}

// SetLevel to change current level
func (logger *Logger) SetLevel(levelLabel string) {
	logger.level = label2level(levelLabel)
}

// Error out
func (logger *Logger) Error(a ...interface{}) {
	if errorLevel > logger.level {
		return
	}
	logger.errorOut.Println(a...)
}

// Errorf out
func (logger *Logger) Errorf(format string, a ...interface{}) {
	logger.Error(fmt.Sprintf(format, a...))
}

// Warning out
func (logger *Logger) Warning(a ...interface{}) {
	if warningLevel > logger.level {
		return
	}
	logger.warningOut.Println(a...)
}

// Warningf out
func (logger *Logger) Warningf(format string, a ...interface{}) {
	logger.Warning(fmt.Sprintf(format, a...))
}

// Info out
func (logger *Logger) Info(a ...interface{}) {
	if infoLevel > logger.level {
		return
	}
	logger.infoOut.Println(a...)
}

// Infof out
func (logger *Logger) Infof(format string, a ...interface{}) {
	logger.Info(fmt.Sprintf(format, a...))
}

// Debug out
func (logger *Logger) Debug(a ...interface{}) {
	if debugLevel > logger.level {
		return
	}
	logger.debugOut.Println(a...)
}

// Debugf out
func (logger *Logger) Debugf(format string, a ...interface{}) {
	logger.Debug(fmt.Sprintf(format, a...))
}

const (
	silentLevel  int = iota // 0
	errorLevel              // 1
	warningLevel            // 2
	infoLevel               // 3
	debugLevel              // 4
)

func label2level(label string) int {
	level := errorLevel
	switch label {
	case "silent":
		level = silentLevel
	case "warning":
		level = warningLevel
	case "info":
		level = infoLevel
	case "debug":
		level = debugLevel
	}

	return level
}
