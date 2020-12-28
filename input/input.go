package input

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/kljensen/snowball"
	"github.com/vudoan2016/ispell/output"
)

type Source struct {
	Title string
	File  string
	Fn    ProcessFn
}

const (
	count   int    = 100
	api_key string = "Token a0ece2037f62563e2f38d2099b31fbc5624b11ab"
)

type owlBotDefinition struct {
	Type       string `json:"type"`
	Definition string `json:"definition"`
	Example    string `json:"Example"`
}

type owlBotDefinitions struct {
	Defs []owlBotDefinition `json:"definitions"`
}

var stopWords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
	"you": {}, "I": {}, "he": {}, "she": {}, "it": {}, "they": {},
}

type ProcessFn func([]rune, *map[string]output.Vocabulary, string) string

func ProcessVocab(line []rune, deck *map[string]output.Vocabulary, title string) string {
	remain := string(line)
	// Indication of the end of an entry which contains a word, a type and a defition.
	if line[len(line)-1] == ';' || line[len(line)-1] == '.' {
		words := strings.Fields(string(line))
		_word := words[0]
		if _, found := (*deck)[_word]; !found {
			_type := words[1]
			(*deck)[_word] =
				output.Vocabulary{
					Word: _word,
					Type: string(_type[:len(_type)-1]),
					Def:  strings.Join(words[2:], " ")}
		}
		remain = ""
	}
	return remain
}

func toSentences(line []rune) ([]string, string) {
	var out []string
	var remain string
	start := 0
	for i := 0; i < len(line); i++ {
		if line[i] == '.' || line[i] == ';' || line[i] == '!' || line[i] == '?' {
			out = append(out, string(line)[start:i+1])
			start = i + 1
		}
	}
	if start < len(line) {
		remain = string(line[start:])
	}
	return out, remain
}

func ProcessBook(line []rune, deck *map[string]output.Vocabulary, title string) string {
	citation := "(" + title + ")"
	punctuation := regexp.MustCompile(`[,.;"?!:]()`)
	if strings.Contains(string(line), "CHAPTER") || strings.Contains(string(line), "PART") {
		return ""
	}
	sentences, remain := toSentences(line)
	for _, s := range sentences {
		words := strings.Fields(s)
		for _, word := range words {
			word = strings.ToLower(word)

			// eliminate punctuation marks
			word = punctuation.ReplaceAllString(word, "")
			// skip stop words
			if _, isStopWord := stopWords[word]; isStopWord {
				continue
			}
			// skip compound words
			if strings.Contains(word, "-") {
				continue
			}
			// skip abbriviations
			if strings.Contains(word, "'") {
				continue
			}
			if _, found := (*deck)[word]; !found {
				(*deck)[word] = output.Vocabulary{Word: word, Usage: s}
			} else {
				(*deck)[word] = output.Vocabulary{Word: word, Type: (*deck)[word].Type, Def: (*deck)[word].Def, Usage: s + " " + citation}
			}
		}
	}
	return remain
}

func cleanText(scanner *bufio.Scanner, deck *map[string]output.Vocabulary, fn ProcessFn, title string) {
	var remain string
	space := regexp.MustCompile(`\s+`)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			// compress whitespaces (spaces, tabs, ...)
			line := space.ReplaceAllString(line, " ")

			remain = fn([]rune(remain+" "+line), deck, title)
		}
	}
}

func readText(source []Source) (map[string]output.Vocabulary, error) {
	deck := make(map[string]output.Vocabulary)

	for _, src := range source {
		file, err := os.Open(src.File)
		if err != nil {
			return deck, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		cleanText(scanner, &deck, src.Fn, src.Title)
		file.Close()
	}
	return deck, nil
}

func getOwlDefinitions(word string) (owlBotDefinitions, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", "https://owlbot.info/api/v4/dictionary/"+word, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", api_key)

	// Send req using http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	var definitions owlBotDefinitions
	responseData, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal([]byte(responseData), &definitions)
		if err != nil {
			log.Println("Failed to unmarshal response for ", word)
		}
	}
	return definitions, err
}

func cleanOutput(word string) string {
	word = strings.ReplaceAll(word, "<b>", " ")
	word = strings.ReplaceAll(word, "</br>", " ")
	word = strings.ReplaceAll(word, "<b>", " ")
	word = strings.ReplaceAll(word, "</b>", " ")
	word = strings.ReplaceAll(word, "<span>", " ")
	word = strings.ReplaceAll(word, "</span>", " ")
	return word
}

func selectWords(vocabs *map[string]output.Vocabulary) ([]output.Vocabulary, error) {
	var i int
	var out []output.Vocabulary
	var wordType, def, usage string

	for _, v := range *vocabs {
		def = v.Def
		wordType = v.Type
		usage = v.Usage
		owlDefs, _ := getOwlDefinitions(v.Word)
		if len(owlDefs.Defs) == 0 {
			// try stemmed word
			stemmed, err := snowball.Stem(v.Word, "english", true)
			if err == nil {
				// is stemmed word valid
				owlDefs, _ = getOwlDefinitions(stemmed)
			}
		}
		if len(owlDefs.Defs) > 0 {
			if len(def) == 0 {
				def = cleanOutput(owlDefs.Defs[0].Definition)
			}
			if len(wordType) == 0 {
				switch owlDefs.Defs[0].Type {
				case "noun":
					wordType = "n"
				case "verb":
					wordType = "v"
				case "adjective":
					wordType = "adj"
				}
			}
			if len(usage) == 0 {
				usage = cleanOutput(owlDefs.Defs[0].Example)
			}
		}
		out = append(out, output.Vocabulary{Word: v.Word, Type: wordType, Def: def, Usage: usage})
		i++
		if i == count {
			break
		}
	}
	return out, nil
}

func Init(source []Source) (map[string]output.Vocabulary, []output.Vocabulary, error) {
	var words []output.Vocabulary
	deck, err := readText(source)
	if err == nil {
		words, err = selectWords(&deck)
	}
	return deck, words, err
}
