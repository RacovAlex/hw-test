package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	testFile, err := os.CreateTemp("./testdata", "testfile-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.Remove(testFile.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	testCases := []struct {
		name   string
		from   string
		to     string
		offset int64
		limit  int64
	}{
		{
			name:   "copy full file",
			from:   "./testdata/input.txt",
			to:     testFile.Name(),
			offset: 0,
			limit:  0,
		},
		{
			name:   "normal case offset=0, limit=1000",
			from:   "./testdata/input.txt",
			to:     testFile.Name(),
			offset: 0,
			limit:  1000,
		},
		{
			name:   "normal case offset=500, limit=0",
			from:   "./testdata/input.txt",
			to:     testFile.Name(),
			offset: 500,
			limit:  0,
		},
		{
			name:   "normal case offset=500, limit=3500",
			from:   "./testdata/input.txt",
			to:     testFile.Name(),
			offset: 500,
			limit:  3500,
		},
		{
			name:   "copy 1 byte",
			from:   "./testdata/input.txt",
			to:     testFile.Name(),
			offset: 500,
			limit:  1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Выполняем копирование
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.NoError(t, err)

			// Получаем данные файлов
			fromFileInfo, err := os.Stat(tc.from)
			require.NoError(t, err)
			toFileInfo, err := os.Stat(tc.to)
			require.NoError(t, err)

			// Проверяем размер файла
			if tc.limit == 0 {
				assert.Equal(t, fromFileInfo.Size()-tc.offset, toFileInfo.Size(), "file sizes should be equal")
			} else {
				assert.Equal(t, tc.limit, toFileInfo.Size())
			}

			// Проверяем содержимое файла
			fromData, err := os.ReadFile(tc.from)
			require.NoError(t, err)
			toData, err := os.ReadFile(tc.to)
			require.NoError(t, err)
			if tc.limit == 0 {
				assert.Equal(t, fromData[tc.offset:], toData, "file contents should be equal")
			} else {
				assert.Equal(t, fromData[tc.offset:(tc.offset+tc.limit)], toData, "file contents should be equal")
			}
			fmt.Println("")
		})
	}

	t.Run("open file error", func(t *testing.T) {
		err := Copy("./testdata/input25.txt", testFile.Name(), 0, 0)
		require.ErrorAs(t, err, new(*os.PathError))
	})

	t.Run("exceeded file boundaries", func(t *testing.T) {
		err := Copy("./testdata/input.txt", testFile.Name(), 8000, 0)
		require.ErrorAs(t, err, &ErrOffsetExceedsFileSize)
	})

	t.Run("file is closed", func(t *testing.T) {
		err := Copy("./testdata/input.txt", testFile.Name(), 0, 0)
		require.NoError(t, err)

		// Пробуем повторно открыть файл для записи
		f, err := os.OpenFile(testFile.Name(), os.O_RDWR, 0o644)
		require.NoError(t, err)
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(f)
	})
}
