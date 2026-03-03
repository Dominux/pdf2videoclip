package internal

import (
	"app/internal/common"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var (
	translitRules = []string{
		"shch", "щ", "sh", "ш", "ch", "ч", "kh", "х", "ts", "ц",
		"yo", "ё", "zh", "ж", "eh", "э", "yu", "ю", "ya", "я",
		"a", "а", "b", "б", "v", "в", "g", "г", "d", "д",
		"e", "е", "z", "з", "i", "и", "j", "й", "k", "к",
		"l", "л", "m", "м", "n", "н", "o", "о", "p", "п",
		"r", "р", "s", "с", "t", "т", "u", "у", "f", "ф",
		"y", "ы", "h", "х", "w", "в", "x", "кс", "q", "к",
	}

	replacer = strings.NewReplacer(translitRules...)
)

type VoiceoverAdapter struct {
	text string
}

func newVoiceoverAdapter(text string) VoiceoverAdapter {
	return VoiceoverAdapter{text}
}

func (s *VoiceoverAdapter) prepareText() {
	// making text one line
	text := regexp.MustCompile(`[^\S ]+`).ReplaceAllString(s.text, " ")

	// lowercasing text
	text = strings.ToLower(text)

	// dividing into "sentences"
	type Entry struct {
		sentence string
		endMark  rune
	}
	sentences := []Entry{}
	{
		runes := []rune{}
		for _, char := range text {
			if slices.Contains([]rune{',', '.', '!', '?'}, char) {
				sentences = append(sentences, Entry{string(runes), char})
				runes = []rune{}
			} else {
				runes = append(runes, char)
			}
		}
	}

	newSentences := []string{}
	for _, entry := range sentences {
		newWords := []string{}

		// dividing into words
		for _, w := range strings.Split(entry.sentence, " ") {
			if w == "" {
				continue
			}

			word := russifyWord(w)
			newWords = append(newWords, word)
		}

		newSentence := strings.Join(newWords, " ") + string(entry.endMark) + " "
		newSentences = append(newSentences, newSentence)
	}

	s.text = strings.Join(newSentences, "")
}

func russifyWord(word string) string {
	r := []rune(word)

	if unicode.Is(unicode.Latin, r[0]) {
		return replacer.Replace(word)
	}

	if unicode.IsDigit(r[0]) {
		num, err := strconv.Atoi(word)
		if err != nil {
			fmt.Println("Error during conversion:", err)
			return word
		}

		return common.NumberToRussian(num)
	}

	return word
}
