package internal

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	sentenceStopTime    float64 = .4
	subsentenceStopTime float64 = .15
)

type renderer struct {
	text string
}

func newRenderer(text string) renderer {
	return renderer{text}
}

type TranscriptedWord struct {
	word  string
	start float64
	end   float64
}

func (r *renderer) render() {
	transcription := r.transcriptText()

	fmt.Println(transcription)
}

func (r *renderer) transcriptText() []TranscriptedWord {
	// 1. splitting text into sentences
	textSentences := regexp.MustCompile(`[.!?]+`).Split(r.text, -1)

	// 2. splitting sentences into sub sentences (by comma)
	subSentenceStopTotal := 0
	sentences := make([][]string, 0, len(textSentences))
	for _, textSentence := range textSentences {
		subsentences := strings.Split(textSentence, ",")
		sentences = append(sentences, subsentences)

		subSentenceStopTotal += len(subsentences) - 1
	}

	// 3. calculating total duration of such time stops
	stopTotalLength := (float64(len(textSentences)) * sentenceStopTime) + (float64(subSentenceStopTotal) * subsentenceStopTime)

	// 4. getting the voice file duration
	duration, err := getDuration(finalFilepath)
	if err != nil {
		panic(err)
	}

	// 5. deleting stops duration from the total duration
	durationWithoutStops := duration - stopTotalLength

	// 6. getting the average duraton of the symbol
	// (without removing stops cuz spaces between words takes more time than other symbols so it's like balancing their duration)
	avgSymbolDuration := durationWithoutStops / float64(len(r.text))

	// 7. finally creating transcript
	transcript := []TranscriptedWord{}
	timeCounter := 0.0
	for _, sentence := range sentences {
		for _, subsentence := range sentence {
			words := strings.Split(subsentence, " ")

			for _, word := range words {
				wordTime := float64(len(word)) * avgSymbolDuration
				nextWordStart := timeCounter + wordTime

				wordTranscription := TranscriptedWord{word: word, start: timeCounter, end: nextWordStart}
				transcript = append(transcript, wordTranscription)

				timeCounter = nextWordStart
			}

			// adding sub sentence duration
			timeCounter += subsentenceStopTime
		}

		// adding sentence duration
		timeCounter += sentenceStopTime
	}

	return transcript
}

func getDuration(fileName string) (float64, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		fileName,
	}

	out, err := exec.Command("ffprobe", args...).Output()
	if err != nil {
		return 0, err
	}

	// 1. Clean the string
	outStr := strings.TrimSpace(string(out))

	// 2. Parse string to float64
	seconds, err := strconv.ParseFloat(outStr, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse duration: %v", err)
	}

	return seconds, nil
}
