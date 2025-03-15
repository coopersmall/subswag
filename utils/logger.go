package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	projectName = "subswag"
)

type ILogger interface {
	Info(ctx context.Context, msg string, metadata map[string]any)
	Debug(ctx context.Context, msg string, metadata map[string]any)
	Warn(ctx context.Context, msg string, metadata map[string]any)
	Error(ctx context.Context, msg string, err error, metadata map[string]any)
}

type stdLogger struct {
	service   string
	stdLogger *slog.Logger
}

type opt func(*stdLogger)

func NewLogger(service string, opts ...opt) ILogger {
	logger := &stdLogger{
		service:   service,
		stdLogger: slog.New(newLogHandler(os.Stdout, &noopLogExporter{})),
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

func WithLoggingHandler(w io.Writer, exporter LogExporter) opt {
	return func(logger *stdLogger) {
		logger.stdLogger = slog.New(newLogHandler(w, exporter))
	}
}

func (l *stdLogger) Info(ctx context.Context, msg string, metadata map[string]any) {
	logLevel(ctx, l.stdLogger, l.service, slog.LevelInfo, msg, nil, metadata, 2)
}
func (l *stdLogger) Debug(ctx context.Context, msg string, metadata map[string]any) {
	logLevel(ctx, l.stdLogger, l.service, slog.LevelDebug, msg, nil, metadata, 2)
}
func (l *stdLogger) Warn(ctx context.Context, msg string, metadata map[string]any) {
	logLevel(ctx, l.stdLogger, l.service, slog.LevelWarn, msg, nil, metadata, 2)
}
func (l *stdLogger) Error(ctx context.Context, msg string, err error, metadata map[string]any) {
	logLevel(ctx, l.stdLogger, l.service, slog.LevelError, msg, err, metadata, 2)
}

func logLevel(
	ctx context.Context,
	Logger *slog.Logger,
	service string,
	level slog.Level,
	msg string,
	err error,
	metadata map[string]any,
	callerDepth int,
) {
	var attrs []slog.Attr
	if metadata != nil {
		attrs = append(attrs, slog.Any("metadata", metadata))
	}

	if level == slog.LevelError {
		if err == nil {
			err = errors.New(msg)
		}
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	source := getSource(callerDepth + 1)
	if source != "" {
		attrs = append(attrs, slog.String("source", source))
	}

	attrs = append(attrs, slog.String("service", service))

	Logger.LogAttrs(ctx, level, msg, attrs...)
}

func getSource(callerDepth int) string {
	pc, file, line, ok := runtime.Caller(callerDepth)
	if !ok {
		return ""
	}

	parts := strings.Split(file, projectName)
	if len(parts) > 1 {
		file = projectName + parts[1]
	} else {
		file = projectName + "/"
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts = strings.Split(funcName, ".")
	funcName = parts[len(parts)-1]

	return fmt.Sprintf("%s:%d:%s", file, line, funcName)
}

type LogExporter interface {
	export(log []byte)
}

type localLogExporter struct {
	file    *os.File
	isEmpty bool
	logs    [][]byte
}

func newLocalExporter() (LogExporter, error) {
	fileName := fmt.Sprintf("%s/%d.log", logDir, time.Now().Unix())
	out, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	exporter := &localLogExporter{
		file:    out,
		isEmpty: true,
		logs:    make([][]byte, 0),
	}

	go exporter.start()

	return exporter, nil
}

func (e *localLogExporter) start() error {
	for {
		select {
		case <-time.After(defaultWriteInterval):
			e.write()
		}
	}
}
func (e *localLogExporter) export(log []byte) {
	e.logs = append(e.logs, log)
}
func (e *localLogExporter) write() error {
	for _, log := range e.logs {
		if !e.isEmpty {
			log = append([]byte(",\n"), log...)
		}

		_, err := e.file.Write(log)
		if err != nil {
			return err
		}

		e.isEmpty = false
	}

	e.logs = make([][]byte, 0)

	return nil
}

type noopLogExporter struct{}

func (e *noopLogExporter) export(log []byte) {}

const (
	logDir               = "./.logs"
	defaultWriteInterval = 10 * time.Second
)

type logHandler struct {
	slog.Handler
	l *log.Logger
	e LogExporter
}

func newLogHandler(writer io.Writer, exporter LogExporter) *logHandler {
	h := &logHandler{
		Handler: slog.NewJSONHandler(writer, &slog.HandlerOptions{
			AddSource: true,
		}),
		l: log.New(writer, "", 0),
		e: exporter,
	}

	return h
}

func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	fields["message"] = r.Message
	fields["level"] = r.Level.String()
	fields["timestamp"] = r.Time.Format(time.RFC3339)

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	h.e.export(b)

	h.l.Println(string(b))

	return nil
}

func GetLogger(name string) ILogger {
	return NewLogger(name)
}
