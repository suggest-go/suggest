package analysis

import (
	snowball "github.com/snowballstem/snowball/go"
	"github.com/suggest-go/suggest/pkg/analysis/en"
	"github.com/suggest-go/suggest/pkg/analysis/ru"
)

type stemmerFunc = func(env *snowball.Env) bool
type stopWordsSet = map[string]struct{}

type stemmerFilter struct {
	stemmer   stemmerFunc
	stopWords stopWordsSet
}

// NewRussianStemmerFilter creates a new stemmer for Russian language.
func NewRussianStemmerFilter() TokenFilter {
	return &stemmerFilter{
		stemmer:   ru.Stem,
		stopWords: stopWordsToSet(ru.StopWords),
	}
}

// NewEnglishStemmerFilter creates a new stemmer for English language.
func NewEnglishStemmerFilter() TokenFilter {
	return &stemmerFilter{
		stemmer:   en.Stem,
		stopWords: stopWordsToSet(en.StopWords),
	}
}

// Filter filters the given list with described behaviour
func (f *stemmerFilter) Filter(list []Token) []Token {
	env := snowball.NewEnv("")
	filtered := []Token{}

	for _, token := range list {
		if _, ok := f.stopWords[token]; ok {
			continue
		}

		env.SetCurrent(token)
		f.stemmer(env)
		filtered = append(filtered, env.Current())
	}

	return filtered
}

func stopWordsToSet(list []string) stopWordsSet {
	set := make(stopWordsSet, len(list))

	for _, word := range list {
		set[word] = struct{}{}
	}

	return set
}
