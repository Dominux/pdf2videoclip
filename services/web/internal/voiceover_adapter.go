package internal

import (
	"app/internal/common"
	"fmt"
	"regexp"
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
	sentences := regexp.MustCompile(`[,.?!]\s*`).Split(text, -1)

	newSentences := []string{}
	for _, s := range sentences {
		newWords := []string{}

		// diving into words
		for _, w := range strings.Split(s, " ") {
			if w == "" {
				continue
			}

			word := russifyWord(w)
			newWords = append(newWords, word)
		}

		newSentence := "<s>" + strings.Join(newWords, " ") + "</s>"
		newSentences = append(newSentences, newSentence)
	}

	s.text = "<speech>" + strings.Join(newSentences, "") + "</speech>"
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
