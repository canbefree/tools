package helper

import (
	"context"
	"log"
)

type LogI interface {
	Printf(string, ...any)
	PrintfWithContext(context.Context, string, ...any)
	Errorf(string, ...any)
	ErrorfWithContext(context.Context, string, ...any)
}

var DefaultLog = &Log{}

type Log struct {
}

func (l *Log) Printf(format string, v ...any) {
	Printf(format, v...)
}

func (l *Log) PrintfWithContext(ctx context.Context, format string, v ...any) {
	PrintfWithContext(ctx, format, v...)
}

func (l *Log) Errorf(format string, v ...any) {
	Errorf(format, v...)
}

func (l *Log) ErrorfWithContext(ctx context.Context, format string, v ...any) {
	ErrorfWithContext(ctx, format, v...)
}

var orginalLog = log.Default()

func Printf(format string, v ...any) {
	orginalLog.Printf(format, v...)
}

func PrintfWithContext(ctx context.Context, format string, v ...any) {
	orginalLog.Printf(format, v...)
}

func Errorf(format string, v ...any) {
	orginalLog.Fatalf(format, v...)
}

func ErrorfWithContext(ctx context.Context, format string, v ...any) {
	orginalLog.Fatalf(format, v...)
}
