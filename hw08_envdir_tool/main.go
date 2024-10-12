package main

import (
	"log/slog"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		slog.Error("Usage: go run main.go <path-to-dir> <command> [args...]")
		os.Exit(1)
	}

	// Чтение переменных окружения из директории
	env, err := ReadDir(args[1])
	if err != nil {
		slog.Error("Error reading directory", slog.String("dir", args[1]), slog.Any("error", err))
		os.Exit(1)
	}

	// Запуск команды с переданными аргументами и переменными окружения
	exitCode := RunCmd(args[2:], env)
	if exitCode != 0 {
		slog.Error("Command execution failed", slog.String("command", args[2]), slog.Int("exitCode", exitCode))
	}

	// Завершаем с кодом возврата команды
	os.Exit(exitCode)
}
