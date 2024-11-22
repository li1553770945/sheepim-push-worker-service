package log

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type TraceLogger struct {
	*kitexlogrus.Logger
}

func NewTraceLogger() *TraceLogger {
	return &TraceLogger{
		Logger: kitexlogrus.NewLogger(),
	}
}

func getCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func (l *TraceLogger) logWithTrace(ctx context.Context, level, msg string) {
	caller := getCallerInfo(4)
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("level", level),
			attribute.String("message", msg),
			attribute.String("caller", caller),
		))
	}
}

func (l *TraceLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "DEBUG", fmt.Sprintf(format, v...))
	l.Logger.CtxDebugf(ctx, format, v...)
}

func (l *TraceLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "INFO", fmt.Sprintf(format, v...))
	l.Logger.CtxInfof(ctx, format, v...)
}

func (l *TraceLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "NOTICE", fmt.Sprintf(format, v...))
	l.Logger.CtxNoticef(ctx, format, v...)
}

func (l *TraceLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "WARN", fmt.Sprintf(format, v...))
	l.Logger.CtxWarnf(ctx, format, v...)
}

func (l *TraceLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "ERROR", fmt.Sprintf(format, v...))
	l.Logger.CtxErrorf(ctx, format, v...)
}

func (l *TraceLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.logWithTrace(ctx, "FATAL", fmt.Sprintf(format, v...))
	l.Logger.CtxFatalf(ctx, format, v...)
}

func InitLog() *TraceLogger {
	logger := NewTraceLogger()
	klog.SetLogger(logger)
	klog.SetLevel(klog.LevelInfo)
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		klog.Fatalf("创建日志文件夹失败：%v", err)
	}

	logFileName := filepath.Join("logs", "log-"+time.Now().Format("2006-01-02-15-04-05")+".log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		klog.Fatalf("failed to open log file: %v", err)
	}
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	klog.SetOutput(multiWriter)

	httpLogFileName := filepath.Join("logs", "http-log-"+time.Now().Format("2006-01-02-15-04-05")+".log")
	httpLogFile, err := os.OpenFile(httpLogFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		klog.Fatalf("failed to open log file: %v", err)
	}
	httpMultiWriter := io.MultiWriter(httpLogFile, os.Stdout)
	hlog.SetOutput(httpMultiWriter)
	return logger
}
