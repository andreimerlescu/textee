package textee

import (
	"errors"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/andreimerlescu/gematria"
)

var (
	ErrEmptyInput    ArgumentError = errors.New("empty input")
	ErrGematriaParse GematriaError = errors.New("unable to parse gematria for value")
	ErrRegexpMissing RegexpError   = errors.New("regexp compile result missing")
	ErrBadParsing    ParseError    = errors.New("failed to parse the string")
)

type ArgumentError error
type GematriaError error
type RegexpError error
type ParseError error
type CleanError error

type Textee struct {
	mu             sync.RWMutex
	Input          string                       `json:"in"`
	Gematria       gematria.Gematria            `json:"gem"`
	Substrings     map[string]*atomic.Int32     `json:"subs"` // map[Substring]*atomic.Int32
	Gematrias      map[string]gematria.Gematria `json:"gems"`
	ScoresEnglish  map[uint64][]string          `json:"sen"`
	ScoresJewish   map[uint64][]string          `json:"sje"`
	ScoresSimple   map[uint64][]string          `json:"ssi"`
	ScoresMystery  map[uint64][]string          `json:"smy"`
	ScoresMajestic map[uint64][]string          `json:"smj"`
	ScoresEights   map[uint64][]string          `json:"sei"`
}

type SubstringQuantity struct {
	Substring string `json:"s"`
	Quantity  int    `json:"q"`
}

type SortedStringQuantities []SubstringQuantity

var regCleanSubstring = regexp.MustCompile(`[^a-zA-Z0-9\s]`)
var regFindSentences = regexp.MustCompile(`(?m)([^.!?]*[.!?])(?:\s|$)`)

// stringToSentenceSlice splits text into sentences, considering abbreviations.
func stringToSentenceSlice(text string) ([]string, error) {
	if regFindSentences == nil {
		return nil, ErrRegexpMissing
	}
	matches := regFindSentences.FindAllString(text, -1)
	for i, match := range matches {
		matches[i] = strings.TrimSpace(match)
	}

	if len(matches) == 0 {
		return []string{text}, nil
	}
	return matches, nil
}

// cleanSubstring returns the string to A-Za-z0-9\s only
func cleanSubstring(word string) (string, error) {
	if regCleanSubstring == nil {
		return "", ErrRegexpMissing
	}
	word = strings.TrimSpace(word)
	return regCleanSubstring.ReplaceAllString(word, ""), nil
}
