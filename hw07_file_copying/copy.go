package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy копирует байты из одного файла в другой
// offset устанавливает отступ для чтения
// limit определяет количество байт для чтения.
func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем файл для чтения
	readFile, readFileInfo, err := OpenFile(fromPath, offset)
	if err != nil {
		return err
	}
	defer CloseFile(readFile, fromPath)

	// Устанавливаем отступ
	_, err = readFile.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("seek in file: %w", err)
	}

	// Открываем файл для записи
	writeFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("open file %s: %w", toPath, err)
	}
	defer CloseFile(writeFile, toPath)

	// Создаем прогресс-бар и оборачиваем его в Reader
	bar := pb.Full.Start64(readFileInfo.Size())
	barReader := bar.NewProxyReader(readFile)
	defer bar.Finish()

	err = copyData(writeFile, barReader, limit)

	return err
}

// OpenFile открывает файл для чтения и возвращает вместе с ним информацию о нем.
func OpenFile(filePath string, offset int64) (*os.File, os.FileInfo, error) {
	readFile, err := os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("open file %s: %w", filePath, err)
	}

	// Получаем данные файла
	fileInfo, err := readFile.Stat()
	if err != nil {
		return nil, nil, ErrUnsupportedFile
	}

	// Проверяем, что отступ не выходит за границы файла
	if fileInfo.Size() < offset {
		return nil, nil, ErrOffsetExceedsFileSize
	}
	return readFile, fileInfo, nil
}

func CloseFile(file *os.File, filePath string) {
	if err := file.Close(); err != nil {
		log.Panicf("close file %s: %s", filePath, err)
	}
}

// copyData копирует байты из одного ридера в другой в зависимости от установленного лимита.
func copyData(writer io.Writer, reader io.Reader, limit int64) error {
	var copier int64
	var err error
	if limit > 0 {
		copier, err = io.CopyN(writer, reader, limit)
	} else {
		copier, err = io.Copy(writer, reader)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("copy file: %w", err)
	}
	fmt.Printf("copied %d bytes\n", copier)
	return nil
}
