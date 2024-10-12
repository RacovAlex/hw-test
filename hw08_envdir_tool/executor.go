package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		slog.Error("No command to run")
		return 1
	}
	// Проверяем, существует ли команда
	bin, err := exec.LookPath(cmd[0])
	if err != nil {
		slog.Error(fmt.Sprintf("Command not found: %v", err))
		return 1
	}

	command := exec.Command(bin, cmd[1:]...)

	// Устанавливаем новые переменные окружения
	command.Env = ChangeEnvVars(os.Environ(), env)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	err = command.Run()
	if err == nil {
		return 0
	}

	// Используем errors.As для обработки *exec.ExitError
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	// Если это не *exec.ExitError, выводим сообщение об ошибке
	slog.Error(fmt.Sprintf("Failed to run command: %v", err))
	return 1
}

// ChangeEnvVars изменяет слайс переменных окружения в соответствии с полученной картой env.
func ChangeEnvVars(envVars []string, env Environment) []string {
	// Преобразуем текущие переменные окружения в карту (map) для удобства поиска
	envMap := make(map[string]string, len(envVars))
	for _, envVar := range envVars {
		str := strings.SplitN(envVar, "=", 2)
		if len(str) != 2 {
			slog.Error("Invalid environment variable", "variable", envVar)
			continue
		}
		envMap[str[0]] = str[1]
	}

	// Изменяем переменные окружения
	for varName, value := range env {
		if value.NeedRemove {
			delete(envMap, varName)
		} else {
			envMap[varName] = value.Value
		}
	}

	// Преобразуем обновленную карту переменных окружения в слайс
	vars := make([]string, 0, len(envMap))
	for varName, varValue := range envMap {
		vars = append(vars, fmt.Sprintf("%s=%s", varName, varValue))
	}
	return vars
}
