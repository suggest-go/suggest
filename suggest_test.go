package suggest

import (
	"os"
	"reflect"
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
		for i := 0; i < 10; i++ {
			service.AddRunTimeIndex(wordsList[i%5], conf)
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

func TestConcurrencyOnDisc(t *testing.T) {
	configFile, err := os.Open("testdata/config.json")
	if err != nil {
		t.Error(err)
	}

	description, err := ReadConfigs(configFile)
	if err != nil {
		t.Error(err)
	}

	wordsList := []string{"Nissan March", "Honda Fitt", "Wolfsvagen", "Tayota Corolla", "Micra Nissan"}
	service := NewService()

	err = service.AddOnDiscIndex(description[0])
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 5; i++ {
			err = service.AddOnDiscIndex(description[0])
			if err != nil {
				t.Error(err)
			}
		}
		wg.Done()
	}()

	expectedValues := [][]string{
		{"NISSAN MARCH"},
		{"HONDA FIT"},
		{},
		{"TOYOTA COROLLA"},
		{"NISSAN MICRA"},
	}

	go func() {
		for i := 0; i < 5; i++ {
			searchConf, _ := NewSearchConfig(wordsList[i], 5, CosineMetric(), 0.7)
			result := service.Suggest(description[0].Name, searchConf)
			actual := make([]string, 0, len(result))
			for _, item := range result {
				actual = append(actual, item.Value)
			}

			if !reflect.DeepEqual(actual, expectedValues[i]) {
				t.Errorf("Test Fail, expected %v, got %v", expectedValues[i], actual)
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
