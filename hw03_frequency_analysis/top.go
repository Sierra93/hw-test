package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

// WordCount структура для хранения слова и его частоты.
type WordCount struct {
	Word  string
	Count int
}

// ByFrequencyAndLex сортировка по убыванию частоты, при равенстве — лексикографически.
type ByFrequencyAndLex []WordCount

func (a ByFrequencyAndLex) Len() int {
	return len(a)
}

func (a ByFrequencyAndLex) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByFrequencyAndLex) Less(i, j int) bool {
	if a[i].Count != a[j].Count {
		return a[i].Count > a[j].Count // по убыванию частоты
	}

	return a[i].Word < a[j].Word // лексикографически
}

/*
Top10 возвращает первые 10 слов после всех сортировок и нужных условий.
*/
func Top10(text string) []string {
	// Используем регулярное выражение для разделения текста на слова,
	// учитывая, что знаки препинания считаются частью слова.
	re := regexp.MustCompile(`[^\s]+`)
	words := re.FindAllString(text, -1)

	freqMap := make(map[string]int)

	// Готовим мапу где ключ это слово, а значение кол-во его повторений.
	for _, word := range words {
		freqMap[word]++
	}

	// Создаем слайс для сортировки.
	var wcSlice []WordCount

	// Заполняем слайс словом и его кол-вом повторений.
	for w, c := range freqMap {
		wcSlice = append(wcSlice, WordCount{w, c})
	}

	// Сортируем слайс слов.
	sort.Sort(ByFrequencyAndLex(wcSlice))

	// Заводим слайс с капасити 10.
	result := make([]string, 0, 10)

	// Берем первые 10 слов.
	for i := 0; i < len(wcSlice) && i < 10; i++ {
		result = append(result, wcSlice[i].Word)
	}

	return result
}
