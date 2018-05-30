package vgram

import (
	"github.com/alldroll/suggest/dictionary"
	"github.com/alldroll/suggest/index"
)

type VGramDictionaryBuilder struct {
	qMin, qMax, threshold uint32
	dictionary            dictionary.Dictionary
}

type VGramDictionary struct{}

func (b *VGramDictionaryBuilder) Build() FrequencyTrie {
	trie := b.buildFrequencyTrie()
	b.pruneFrequencyTrie(trie, b.threshold)
	return trie
}

func (b *VGramDictionaryBuilder) buildFrequencyTrie() FrequencyTrie {
	freqTrie := NewFrequencyTrie(b.qMin)
	iter := b.dictionary.Iterator()

	for {
		_, word := iter.GetPair()
		if len(word) == 0 {
			continue
		}

		b.addWord(freqTrie, word)

		if !iter.Next() {
			break
		}
	}

	return freqTrie
}

func (b *VGramDictionaryBuilder) addWord(freqTrie FrequencyTrie, word string) {
	for _, gram := range index.SplitIntoNGrams(word, int(b.qMax)) {
		freqTrie.Add(gram)
	}

	runes := []rune(word)
	lenC := len(runes)

	for q := b.qMax - 1; q >= b.qMin; q-- {
		p := lenC - int(q)
		if p < 0 {
			continue
		}

		substr := string(runes[p:])

		for _, gram := range index.SplitIntoNGrams(substr, int(q)) {
			freqTrie.Add(gram)
		}
	}
}

func (b *VGramDictionaryBuilder) pruneFrequencyTrie(freqTrie FrequencyTrie, threshold uint32) {
	freqTrie.Prune(threshold)
}
