package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	wordsMap := map[string]int{}
	words := strings.Fields(s)

	for _, w := range words {
		wordsMap[w]++
	}

	if len(words) == 0 {
		return nil
	}

	keys := make([]string, 0, len(wordsMap))

	for key := range wordsMap {
		keys = append(keys, key)
	}

	n := 10
	if len(keys) < 10 {
		n = len(keys)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		if wordsMap[keys[i]] == wordsMap[keys[j]] {
			return keys[i] < keys[j]
		}

		return wordsMap[keys[i]] > wordsMap[keys[j]]
	})

	return keys[:n]
}
