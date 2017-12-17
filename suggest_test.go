package suggest

import (
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
	alphabet := NewCompositeAlphabet([]Alphabet{
		NewEnglishAlphabet(),
		NewSimpleAlphabet([]rune{'$'}),
	})

	wordsList := []string{"abc", "test2", "test3", "test4", "teta"}
	dictionary := NewInMemoryDictionary(wordsList)
	conf, _ := NewIndexConfig(3, dictionary, alphabet, "$", "$")
	service := NewService()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 5; i++ {
			service.AddRunTimeIndex(wordsList[i], conf)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 5; i++ {
			searchConf, _ := NewSearchConfig(wordsList[i], 5, CosineMetric(), 0.7)
			service.Suggest(wordsList[i], searchConf)
		}

		wg.Done()
	}()

	wg.Wait()
}
