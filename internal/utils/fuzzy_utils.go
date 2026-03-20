package utils

import (
	"strings"
	"unicode/utf8"
)

func FuzzyFindWord(needle, haystack string) string {
	needle = strings.TrimSpace(needle)

	if needle == "" {
		return ""
	}

	rawFields := strings.Fields(haystack)
	haystackWords := make([]string, len(rawFields))
	for i, word := range rawFields {
		haystackWords[i] = strings.ToLower(cleanWordBoundary(word))
	}

	words := findWords(needle, haystackWords)
	if len(words) > 0 {
		return words
	}

	for _, word := range haystackWords {
		if strings.HasPrefix(word, needle) {
			return word
		}
	}

	return findAcronym(needle, haystackWords)
}

func findWords(needle string, haystackWords []string) string {
	needleWords := strings.Fields(needle)
	if len(needleWords) < 2 {
		return ""
	}

	foundWords := make([]string, 0, len(needleWords))
	lastHaystackIdx := -1

	for _, needleWord := range needleWords {
		matched := false
		for hIdx := lastHaystackIdx + 1; hIdx < len(haystackWords); hIdx++ {
			word := haystackWords[hIdx]
			if strings.HasPrefix(word, needleWord) {
				foundWords = append(foundWords, word)
				lastHaystackIdx = hIdx
				matched = true
				break
			}
		}
		if !matched {
			return ""
		}
	}

	return strings.Join(foundWords, " ")
}

func findAcronym(needle string, words []string) string {
	if utf8.RuneCountInString(needle) > len(words) {
		return ""
	}

	foundWords := make([]string, 0, utf8.RuneCountInString(needle))
	lastIdx := -1

	for _, nr := range needle {
		matched := false
		for hIdx := lastIdx + 1; hIdx < len(words); hIdx++ {
			hWord := words[hIdx]
			hRune, _ := utf8.DecodeRuneInString(hWord)
			if nr == hRune {
				foundWords = append(foundWords, hWord)
				lastIdx = hIdx
				matched = true
				break
			}
		}
		if !matched {
			return ""
		}
	}

	if len(foundWords) > 0 {
		return strings.Join(foundWords, " ")
	}

	return ""
}

func cleanWordBoundary(word string) string {
	return strings.Trim(word, `'"()[]`)
}
