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
	// Isolates punctuation and groups like ... or !? so they can be processed independently
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
		// Handle basic modifiers
		if words[i] == "(hex)" {
			if i > 0 {
				words[i-1] = hexToDecimal(words[i-1])
			}
			words = append(words[:i], words[i+1:]...)
			i--
		} else if words[i] == "(bin)" {
			if i > 0 {
				words[i-1] = binToDecimal(words[i-1])
			}
			words = append(words[:i], words[i+1:]...)
			i--
		} else if words[i] == "(up)" || words[i] == "(low)" || words[i] == "(cap)" {
			mod := strings.Trim(words[i], "()")
			if i > 0 {
				applyModification(words, i-1, 1, mod)
			}
			words = append(words[:i], words[i+1:]...)
			i--
		} else if i+1 < len(words) && (words[i] == "(up," || words[i] == "(low," || words[i] == "(cap,") {
			// Handle parameterized modifiers like (up, 2)
			mod := strings.Trim(words[i], "(,")
			countStr := strings.Trim(words[i+1], ")")
			count, err := strconv.Atoi(countStr)

			if err == nil && i > 0 {
				applyModification(words, i-1, count, mod)
			}
			// Remove both the tag and the number from the slice
			words = append(words[:i], words[i+2:]...)
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
		if isPunctuation(word) && len(result) > 0 {
			// Attach punctuation strictly to the previous word
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
		if word == "'" {
			if !isOpen {
				isOpen = true
				result = append(result, word) // Add as a standalone opening quote temporarily
			} else {
				if len(result) > 0 {
					result[len(result)-1] += "'" // Attach closing quote to the previous word
				}
				isOpen = false
			}
		} else {
			// If a quote is currently open and waiting for its first word
			if isOpen && len(result) > 0 && result[len(result)-1] == "'" {
				result[len(result)-1] = "'" + word
			} else {
				result = append(result, word)
			}
		}
	}
	return result
}

func fixArticles(words []string) []string {
	for i := 0; i < len(words)-1; i++ {
		if words[i] == "a" || words[i] == "A" {
			nextWord := words[i+1]
			if len(nextWord) > 0 {
				// Strip any leading quotes to accurately check the first letter of the actual word
				firstChar := rune(nextWord[0])
				if firstChar == '\'' && len(nextWord) > 1 {
					firstChar = rune(nextWord[1])
				}

				firstChar = rune(strings.ToLower(string(firstChar))[0])
				if isVowelOrH(firstChar) {
					if words[i] == "a" {
						words[i] = "an"
					} else {
						words[i] = "An"
					}
				}
			}
		}
	}
	return words
}

// ── Helper Utilities ──────────────────────────────────────────

func hexToDecimal(s string) string {
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return s
	}
	return strconv.Itoa(int(n))
}

func binToDecimal(s string) string {
	n, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return s
	}
	return strconv.Itoa(int(n))
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func isPunctuation(s string) bool {
	return s == "." || s == "," || s == "!" || s == "?" || s == ":" || s == ";" || s == "..." || s == "!?"
}

func isVowelOrH(r rune) bool {
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'h'
}
