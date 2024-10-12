package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Читаем директорию с файлами, содержащими значения переменных окружения
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	env := make(Environment, len(files))

	// Заполняем Environment значениями EnvValue (имя переменной: значение, флаг удаления)
	for _, file := range files {
		fileName := file.Name()
		if strings.Contains(fileName, "=") {
			// slog.Error("file %s contains invalid characters =", fileName)
			continue
		}
		filePath := filepath.Join(dir, file.Name())
		openFile, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("open file %s: %w", filePath, err)
		}

		// Читаем первую строку файла, обрабатываем и присваиваем EnvValue ее в качестве значения
		var envVar EnvValue
		scanner := bufio.NewScanner(openFile)
		if scanner.Scan() {
			firstLine := scanner.Text()
			envVar.Value = strings.ReplaceAll(strings.TrimRight(firstLine, " \t"), "\x00", "\n")
			env[fileName] = envVar
		} else {
			envVar.NeedRemove = true
			env[fileName] = envVar
		}
		if err := scanner.Err(); err != nil {
			slog.Error("scan file %s: %w", filePath, err)
		}

		// Закрываем файл до перехода на следующую итерацию
		if err := openFile.Close(); err != nil {
			return nil, fmt.Errorf("close file %s: %w", filePath, err)
		}
	}
	return env, nil
}
