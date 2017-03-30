package suggest

import (
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
	service := NewSuggestService()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		wordsList := []string{"abc", "test2", "test3", "test4", "teta"}
		for i := 0; i < 5; i++ {
			service.AddDictionary(wordsList[i], wordsList, &Config{3, &JaccardDistance{3}, 3, "test"})
		}
		wg.Done()
	}()

	go func() {
		wordsList := []string{"abc", "test2", "test3", "test4", "tetsa"}
		for i := 0; i < 5; i++ {
			service.Suggest(wordsList[i], wordsList[i])
		}

		wg.Done()
	}()

	wg.Wait()
}
