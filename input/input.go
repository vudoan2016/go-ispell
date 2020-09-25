package input

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

const api_key string = "Token a0ece2037f62563e2f38d2099b31fbc5624b11ab"

type Vocabulary struct {
	Word string
	Type string
	Def  string
}

var Deck []Vocabulary

func cleanText(scanner *bufio.Scanner) []Vocabulary {
	var entries []Vocabulary
	var text []string
	space := regexp.MustCompile(`\s+`)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			// compress whitespaces (spaces, tabs, ...)
			line := space.ReplaceAllString(line, " ")

			// split into words
			text = append(text, strings.Fields(line)...)

			// Indication of the end of an entry which contains a word, a type and a defition.
			if line[len(line)-1] == ';' || line[len(line)-1] == '.' {
				_type := text[1]
				_word := text[0]
				entries = append(entries,
					Vocabulary{
						Word: _word,
						Type: string(_type[:len(_type)-1]),
						Def:  strings.Join(text[2:], " ")})
				text = nil
			}
		}
	}
	return entries
}

func splitSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == '.' || data[i] == ';' {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	// There is one final token to be delivered, which may be the empty string.
	// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
	// but does not trigger an error to be returned from Scan itself.
	return 0, data, bufio.ErrFinalToken
}

func readText(books []string) error {
	for _, book := range books {
		file, err := os.Open(book)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		Deck = append(Deck, cleanText(scanner)...)
		file.Close()
	}
	return nil
}

func Init(books []string) error {
	err := readText(books)
	return err
}
