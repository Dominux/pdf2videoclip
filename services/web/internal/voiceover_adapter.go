package internal

import (
	"app/internal/common"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

const finalFilepath = "voice.wav"

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

func (s *VoiceoverAdapter) generateVoice() {
	chunks := s.prepareText()

	// case for a single chunk
	if len(chunks) == 1 {
		if err := s.makeRequest(chunks[0], finalFilepath); err != nil {
			panic(err)
		}
		return
	}

	for i, chunk := range chunks {
		filename := strconv.Itoa(i) + ".wav"
		if err := s.makeRequest(chunk, filename); err != nil {
			panic(err)
		}
	}

	// The full bash command as a string
	script := `ffmpeg -f concat -safe 0 -i <(for f in *.wav; do echo "file '$PWD/$f'"; done) -c copy ` + finalFilepath

	cmd := exec.Command("bash", "-c", script)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func (s *VoiceoverAdapter) makeRequest(text string, filepath string) error {
	apiUrl := "http://localhost:8000/generate"

	// 1. Create url.Values and add your params
	params := url.Values{}
	params.Add("text", text)
	params.Add("speaker", "aidar")
	params.Add("pitch", "40")
	params.Add("rate", "55")

	// 2. Parse the base URL and attach the encoded query
	fullUrl, _ := url.Parse(apiUrl)
	fullUrl.RawQuery = params.Encode()

	logger := log.Default()
	logger.Println("Before Voiceover request")

	resp, err := http.Get(fullUrl.String())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	logger.Println("After Voiceover request")

	// 2. Check if the server actually sent a file (200 OK)
	if resp.StatusCode != http.StatusOK {
		panic("failed to download file: " + resp.Status)
	}

	// 2. Create the local file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	// Ensure the file is closed after the copy operation
	defer out.Close()

	// 3. Use io.Copy to stream data from the response body to the file
	// The response body (resp.Body) is an io.Reader, and the file (out) is an io.Writer.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *VoiceoverAdapter) prepareText() []string {
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

	// converting to chunks cuz of the voiceover limitation
	chunks := []string{}
	newChunk := ""
	for _, sentence := range newSentences {
		if len(newChunk)+len(sentence) >= 900 {
			chunks = append(chunks, newChunk)
			newChunk = ""
		}

		newChunk += sentence
	}
	if newChunk != "" {
		chunks = append(chunks, newChunk)
	}

	return chunks
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
