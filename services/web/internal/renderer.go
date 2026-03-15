package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	r._render(transcription)
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

func (r *renderer) _render(transcript []TranscriptedWord) {
	// 2. Generate .ass subtitle file (better for word-level positioning)
	assFile := "temp_subs.ass"
	f, _ := os.Create(assFile)
	defer os.Remove(assFile)
	defer f.Close()

	// Minimal ASS header
	header := "[Script Info]\n" +
		"ScriptType: v4.00+\n" +
		"[V4+ Styles]\n" +
		"Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\n" +
		"Style: Default,Arial,24,&H00FFFFFF,&H00000000,&H00000000,&H80000000,1,0,0,0,100,100,0,0,3,1,0,2,10,10,10,1\n" +
		"[Events]\n" +
		"Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n"

	f.WriteString(header)

	for _, w := range transcript {
		start := formatTimestamp(w.start)
		end := formatTimestamp(w.end)
		// Format: Dialogue: 0,0:00:00.00,0:00:00.00,Default,,0,0,0,,Word
		line := fmt.Sprintf("Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n", start, end, w.word)
		f.WriteString(line)
	}

	// 3. Execute FFmpeg command
	// -i: input video, -i: input audio, -vf: subtitle filter, -c:v: video codec, -c:a: audio codec
	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", "input.mp4",
		"-i", finalFilepath,
		"-vf", fmt.Sprintf("subtitles=%s", assFile),
		"-c:v", "libx264",
		"-c:a", "aac",
		"-map", "0:v:0", // Use video from first input
		"-map", "1:a:0", // Use audio from second input
		"-shortest", // End when shortest stream ends
		"output.mp4",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Define the command: rm /path/to/dir/*.wav
	// Replace "/path/to/dir/" with your target directory
	cmd = exec.Command("bash", "-c", "rm ./*.wav")

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error deleting files: %v\n", err)
		return
	}

	fmt.Println("All .wav files deleted successfully.")
}

// formatTimestamp converts seconds to ASS format H:MM:SS.CC
func formatTimestamp(seconds float64) string {
	d := time.Duration(seconds * float64(time.Second))
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := float64(d) / float64(time.Second)
	return fmt.Sprintf("%d:%02d:%05.2f", h, m, s)
}
