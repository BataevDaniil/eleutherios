package logging

import (
	"io"
	"log/slog"
	"os"
)

var (
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	logFile *os.File
	logPath string
)

func Configure(path string) error {
	writer := io.Writer(os.Stdout)
	logPath = path
	if path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		logFile = file
		writer = io.MultiWriter(os.Stdout, file)
	}

	logger = slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return nil
}

func Logger() *slog.Logger {
	return logger
}

func LogFilePath() string {
	return logPath
}

func Close() error {
	if logFile == nil {
		return nil
	}
	err := logFile.Close()
	logFile = nil
	return err
}
