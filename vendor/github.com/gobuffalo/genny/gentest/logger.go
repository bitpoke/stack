package gentest

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/gobuffalo/genny"
	"github.com/markbates/safe"
)

const (
	DEBUG string = "DEBU"
	INFO         = "INFO"
	WARN         = "WARN"
	ERROR        = "ERRO"
	FATAL        = "FATA"
	PANIC        = "PANI"
	PRINT        = "PRIN"
)

var _ genny.Logger = &Logger{}

func NewLogger() *Logger {
	l := &Logger{
		Stream: &bytes.Buffer{},
		Log:    map[string][]string{},
		moot:   &sync.Mutex{},
	}
	return l
}

type Logger struct {
	Stream  *bytes.Buffer
	Log     map[string][]string
	PrintFn func(...interface{})
	CloseFn func() error
	moot    *sync.Mutex
}

// Close ...
func (l *Logger) Close() error {
	if l.CloseFn == nil {
		return nil
	}
	return l.CloseFn()
}

func (l *Logger) logf(lvl string, s string, args ...interface{}) {
	l.log(lvl, fmt.Sprintf(s, args...))
}

func (l *Logger) log(lvl string, args ...interface{}) {
	l.moot.Lock()
	m := l.Log[lvl]
	s := fmt.Sprint(args...)
	m = append(m, s)
	l.Stream.WriteString(fmt.Sprintf("[%s] %s\n", lvl, s))
	l.Log[lvl] = m
	l.moot.Unlock()
	if l.PrintFn != nil {
		safe.Run(func() {
			l.PrintFn(args...)
		})
	}
}

func (l *Logger) Debugf(s string, args ...interface{}) {
	l.logf(DEBUG, s, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(DEBUG, args...)
}

func (l *Logger) Infof(s string, args ...interface{}) {
	l.logf(INFO, s, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log(INFO, args...)
}

func (l *Logger) Printf(s string, args ...interface{}) {
	l.logf(PRINT, s, args...)
}

func (l *Logger) Print(args ...interface{}) {
	l.log(PRINT, args...)
}

func (l *Logger) Warnf(s string, args ...interface{}) {
	l.logf(WARN, s, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(WARN, args...)
}

func (l *Logger) Errorf(s string, args ...interface{}) {
	l.logf(ERROR, s, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(ERROR, args...)
}

func (l *Logger) Fatalf(s string, args ...interface{}) {
	l.logf(FATAL, s, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log(FATAL, args...)
}

func (l *Logger) Panicf(s string, args ...interface{}) {
	l.logf(PANIC, s, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.log(PANIC, args...)
}
