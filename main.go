package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input.txt> <output.txt>")
		return
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	final := processText(string(content))
	err = os.WriteFile(os.Args[2], []byte(final), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
}

// processText acts as the orchestrator, passing text through sequential logic blocks
func processText(text string) string {
	// Step 1: Pre-process to separate punctuation from words
	text = preProcess(text)
	words := strings.Fields(text)

	// Step 2: Apply logic modifications
	words = processModifiers(words)
	words = fixPunctuation(words)
	words = fixQuotes(words)
	words = fixArticles(words)

	return strings.Join(words, " ")
}

func preProcess(text string) string {
	// Isolates punctuation and groups like ... or !?
	replacer := strings.NewReplacer(
		"...", " ... ",
		"!?", " !? ",
		".", " . ",
		",", " , ",
		"!", " ! ",
		"?", " ? ",
		":", " : ",
		";", " ; ",
		"'", " ' ",
	)
	return replacer.Replace(text)
}

func processModifiers(words []string) []string {
	for i := 0; i < len(words); i++ {
		word := words[i]
		cleanWord := strings.Trim(word, "(),")

		switch cleanWord {
		case "hex":
			if i > 0 {
				words[i-1] = hexToDecimal(words[i-1])
				words = append(words[:i], words[i+1:]...)
				i--
			}
		case "bin":
			if i > 0 {
				words[i-1] = binToDecimal(words[i-1])
				words = append(words[:i], words[i+1:]...)
				i--
			}
		case "up", "low", "cap":
			count := 1
			// Check for (modifier, N) pattern
			if i+1 < len(words) {
				if n, err := strconv.Atoi(strings.Trim(words[i+1], ")")); err == nil {
					count = n
					words = append(words[:i+1], words[i+2:]...) // Remove the N
				}
			}
			applyModification(words, i-1, count, cleanWord)
			words = append(words[:i], words[i+1:]...) // Remove the modifier tag
			i--
		}
	}
	return words
}

func applyModification(words []string, startIdx, count int, mod string) {
	for j := 0; j < count && startIdx-j >= 0; j++ {
		idx := startIdx - j
		switch mod {
		case "up":
			words[idx] = strings.ToUpper(words[idx])
		case "low":
			words[idx] = strings.ToLower(words[idx])
		case "cap":
			words[idx] = capitalize(words[idx])
		}
	}
}

func fixPunctuation(words []string) []string {
	var result []string
	for _, word := range words {
		if isPunctuation(word[0]) && len(result) > 0 {
			// Attach punctuation to the previous word
			result[len(result)-1] += word
		} else {
			result = append(result, word)
		}
	}
	return result
}

func fixQuotes(words []string) []string {
	var result []string
	isOpen := false
	for _, word := range words {
		if word ==
