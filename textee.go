package textee

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/andreimerlescu/gematria"
)

func NewTextee(opts ...string) (*Textee, error) {
	if opts == nil {
		return nil, ErrEmptyInput
	}

	// ensure regexp clean string compiles
	var cleanRegErr RegexpError
	regCleanSubstring, cleanRegErr = regexp.Compile(`[^a-zA-Z0-9\s]`)
	if cleanRegErr != nil {
		return nil, errors.Join(ErrRegexpMissing, cleanRegErr)
	}

	// ensure regexp find sentences compiles
	var findRegErr RegexpError
	regFindSentences, findRegErr = regexp.Compile(`(?m)([^.!?]*[.!?])(?:\s|$)`)
	if findRegErr != nil {
		return nil, errors.Join(ErrRegexpMissing, findRegErr)
	}

	input := strings.Join(opts, " ")
	gem, err := gematria.NewGematria(input)
	if err != nil {
		return nil, err
	}
	tt := &Textee{
		Input:          input,
		Gematria:       gem,
		Substrings:     make(map[string]*atomic.Int32),
		Gematrias:      make(map[string]gematria.Gematria),
		ScoresEnglish:  make(map[uint64][]string),
		ScoresJewish:   make(map[uint64][]string),
		ScoresSimple:   make(map[uint64][]string),
		ScoresMystery:  make(map[uint64][]string),
		ScoresEights:   make(map[uint64][]string),
		ScoresMajestic: make(map[uint64][]string),
	}
	payload := strings.Join(opts, " ")
	tt, err = tt.ParseString(payload)
	if err != nil {
		return nil, errors.Join(ErrBadParsing, err)
	}
	return tt, nil
}

func (tt *Textee) ParseString(input string) (*Textee, error) {
	sentences, err := stringToSentenceSlice(input)
	if err != nil {
		return nil, errors.Join(ErrBadParsing, err)
	}

	tt.mu.Lock()
	tt.Substrings = make(map[string]*atomic.Int32)
	tt.mu.Unlock()

	var errs []CleanError
	var wg sync.WaitGroup
	for _, sentence := range sentences {
		wg.Add(1)
		go func(sentence string) {
			defer wg.Done()
			words := strings.Fields(sentence)

			for i := 0; i < len(words); i++ {
				for j := i + 1; j <= i+3 && j <= len(words); j++ {
					substring := strings.Join(words[i:j], " ")
					cleanedSubstring, cleanErr := cleanSubstring(substring)
					if cleanErr != nil {
						errs = append(errs, cleanErr)
						continue
					}
					cleanedSubstring = strings.ToLower(cleanedSubstring)

					if cleanedSubstring != "" {
						tt.mu.Lock()
						if _, ok := tt.Substrings[cleanedSubstring]; !ok {
							tt.Substrings[cleanedSubstring] = new(atomic.Int32)
						}
						tt.Substrings[cleanedSubstring].Add(1)
						tt.mu.Unlock()
					}
				}
			}
		}(sentence)
	}
	wg.Wait()
	if len(errs) > 0 {
		for _, e := range errs {
			err = errors.Join(err, e)
		}
		return nil, err
	}
	return tt, nil
}

func (tt *Textee) String() string {
	if len(tt.Substrings) == 0 {
		return ""
	}
	hasGematria := len(tt.ScoresEnglish) > 0 || len(tt.ScoresJewish) > 0 || len(tt.ScoresSimple) > 0
	var output strings.Builder
	for _, data := range tt.SortedSubstrings() {
		if hasGematria := hasGematria; hasGematria {
			output.WriteString(fmt.Sprintf("\"%v\": %d [English %d] [Jewish %d] [Simple %d] [Mystery %d] [Majestic %d] [Eights %d]\n",
				data.Substring, data.Quantity,
				tt.Gematrias[data.Substring].English,
				tt.Gematrias[data.Substring].Jewish,
				tt.Gematrias[data.Substring].Simple,
				tt.Gematrias[data.Substring].Mystery,
				tt.Gematrias[data.Substring].Majestic,
				tt.Gematrias[data.Substring].Eights))
		} else {
			output.WriteString(fmt.Sprintf("\"%v\": %d\n", data.Substring, data.Quantity))
		}
	}
	return output.String()
}

func (tt *Textee) SortedSubstrings() SortedStringQuantities {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	var sortedQuantities SortedStringQuantities

	for k, v := range tt.Substrings {
		quantity := int(v.Load())
		sortedQuantities = append(sortedQuantities, SubstringQuantity{Substring: k, Quantity: quantity})
	}
	sort.Sort(sortedQuantities)

	return sortedQuantities
}

func (tt *Textee) CalculateGematria() (*Textee, error) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	substrings := tt.Substrings
	englishResults := make(map[uint64][]string)
	jewishResults := make(map[uint64][]string)
	simpleResults := make(map[uint64][]string)
	mysteryResults := make(map[uint64][]string)
	majesticResults := make(map[uint64][]string)
	eightsResults := make(map[uint64][]string)
	errorCounter := atomic.Int32{}
	errs := make([]error, 0)
	for substring, _ := range substrings {
		gemscore, err := gematria.NewGematria(substring)
		if err != nil {
			errorCounter.Add(1)
			errs = append(errs, errors.Join(ErrGematriaParse, err))
			continue
		}
		englishResults[gemscore.English] = append(englishResults[gemscore.English], substring)
		jewishResults[gemscore.Jewish] = append(jewishResults[gemscore.Jewish], substring)
		simpleResults[gemscore.Simple] = append(simpleResults[gemscore.Simple], substring)
		mysteryResults[gemscore.Mystery] = append(mysteryResults[gemscore.Mystery], substring)
		majesticResults[gemscore.Majestic] = append(majesticResults[gemscore.Majestic], substring)
		eightsResults[gemscore.Eights] = append(eightsResults[gemscore.Eights], substring)
		if tt.Gematrias == nil {
			tt.Gematrias = make(map[string]gematria.Gematria)
		}
		tt.Gematrias[substring] = gemscore
	}
	if errorCounter.Load() > 0 {
		return nil, errors.Join(errs...)
	}
	substrings = nil
	tt.ScoresEnglish = englishResults
	tt.ScoresJewish = jewishResults
	tt.ScoresSimple = simpleResults
	tt.ScoresMystery = mysteryResults
	tt.ScoresMajestic = majesticResults
	tt.ScoresEights = eightsResults
	englishResults = nil
	jewishResults = nil
	simpleResults = nil
	mysteryResults = nil
	majesticResults = nil
	eightsResults = nil
	return tt, nil
}
