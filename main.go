package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run. <input.txt> <output.txt>")
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

func processText(text string) string {
	text = preProcess(text)
	words := strings.Fields(text)

	words = processModifiers(words)
	words = fixPunctuation(words)
	words = fixQuotes(words)
	words = fixArticles(words)

	return strings.Join(words, " ")
}

func preProcess(text string) string {
	modifiers := []string{
		"(up,", "(low,", "(cap,",
		"(up)", "(low)", "(cap)",
		"(hex)", "(bin)",
	}
	placeholders := make(map[string]string)
	for idx, mod := range modifiers {
		ph := fmt.Sprintf("__MOD_%d__", idx)
		placeholders[ph] = mod
		text = strings.ReplaceAll(text, mod, ph)
	}

	var result strings.Builder
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		if ch == '.' && i+2 < len(runes) && runes[i+1] == '.' && runes[i+2] == '.' {
			result.WriteString("... ")
			i += 2
			continue
		}
		if ch == '!' && i+1 < len(runes) && runes[i+1] == '?' {
			result.WriteString("!? ")
			i++
			continue
		}

		if ch == '\'' {
			if i > 0 && i < len(runes)-1 && isLetter(runes[i-1]) && isLetter(runes[i+1]) {
				result.WriteRune(ch)
			} else {
				result.WriteString(" ' ")
			}
			continue
		}

		if ch == '(' || ch == '"' || ch == '—' {
			result.WriteString(" " + string(ch) + " ")
			continue
		}
		// التعديل هنا: قفل القوس ما يتفصلش لو لازق في كلمة
		if ch == ')' {
			if i > 0 && isLetter(runes[i-1]) {
				result.WriteRune(ch)
			} else {
				result.WriteString(" " + string(ch) + " ")
			}
			continue
		}

		if isPunctuationRune(ch) {
			result.WriteString(" " + string(ch) + " ")
			continue
		}

		result.WriteRune(ch)
	}

	output := result.String()
	for ph, mod := range placeholders {
		output = strings.ReplaceAll(output, ph, mod)
	}

	return output
}

func processModifiers(words []string) []string {
	for i := 0; i < len(words); i++ {
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
			mod := strings.Trim(words[i], "(,")
			countStr := strings.Trim(words[i+1], ")")
			count, err := strconv.Atoi(countStr)

			if err == nil && i > 0 {
				applyModification(words, i-1, count, mod)
			}
			words = append(words[:i], words[i+2:]...)
			i--
		}
	}
	return words
}

func applyModification(words []string, startIdx, count int, mod string) {
	applied := 0
	idx := startIdx
	for applied < count && idx >= 0 {
		if isPunctuation(words[idx]) {
			idx--
			continue
		}
		switch mod {
		case "up":
			words[idx] = strings.ToUpper(words[idx])
		case "low":
			words[idx] = strings.ToLower(words[idx])
		case "cap":
			words[idx] = capitalize(words[idx])
		}
		applied++
		idx--
	}
}

func fixPunctuation(words []string) []string {
	var result []string
	for _, word := range words {
		if isPunctuation(word) && len(result) > 0 {
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
				result = append(result, word)
			} else {
				if len(result) > 0 {
					result[len(result)-1] += "'"
				}
				isOpen = false
			}
		} else {
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
				firstChar := rune(nextWord[0])
				if firstChar == '\'' && len(nextWord) > 1 {
					firstChar = rune(nextWord[1])
				}

				firstChar = rune(strings.ToLower(string(firstChar))[0])
				if isVowel(firstChar) {
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

func isPunctuationRune(r rune) bool {
	return r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';'
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isVowel(r rune) bool {
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u'
}
func removeStrayParentheses(words []string) []string {
	var result []string
	for i, word := range words {
		// extra step to avoid (())
		if word == ")" {
			if i > 0 && !strings.HasSuffix(words[i-1], ")") && !isPunctuation(words[i-1]) {
				continue 
			}
		}
		result = append(result, word)
	}
	return result
}
