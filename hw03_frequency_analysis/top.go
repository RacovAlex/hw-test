package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

// WordFrequency используется для сортировки слов по частотности.
type WordFrequency struct {
	Word  string
	Count int
}

func normalizeWord(word string) string {
	// Приводим к нижнему регистру.
	word = strings.ToLower(word)
	// Удаляем знаки препинания с краев слов.
	word = strings.TrimFunc(word, unicode.IsPunct)
	return word
}

func Top10(s string) []string {
	if s == "" {
		return nil
	}

	// Считаем частотность каждого слова в строке и записываем ее в слайс WordFrequency.
	m := make(map[string]int)
	for _, str := range strings.Fields(s) {
		str = normalizeWord(str)
		if str != "" {
			m[str]++
		}
	}
	wordFreq := make([]WordFrequency, 0, len(m))
	for k, v := range m {
		wordFreq = append(wordFreq, WordFrequency{k, v})
	}

	// Сортируем слайс структур WordFrequency по частотности,
	// а при ее совпадении - лексиграфически.
	sort.Slice(wordFreq, func(i, j int) bool {
		if wordFreq[i].Count == wordFreq[j].Count {
			return wordFreq[i].Word < wordFreq[j].Word
		}
		return wordFreq[i].Count > wordFreq[j].Count
	})

	// Записываем 10 самых частых слов в слайс string.
	index := 10
	if len(wordFreq) < index {
		index = len(wordFreq)
	}
	topWords := make([]string, 0, index)
	for _, k := range wordFreq[:index] {
		topWords = append(topWords, k.Word)
	}

	return topWords
}
