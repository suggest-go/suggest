package suggest

import (
	"reflect"
	"sync"
	"testing"

	"github.com/suggest-go/suggest/pkg/metric"
)

func TestConcurrencyOnDisc(t *testing.T) {
	testConcurrency(t, DiscDriver)
}

func TestConcurrencyRAM(t *testing.T) {
	testConcurrency(t, RAMDriver)
}

func testConcurrency(t *testing.T, driver Driver) {
	descriptions, err := ReadConfigs("testdata/config.json")

	if err != nil {
		t.Error(err)
	}

	description := descriptions[0]
	description.Driver = driver
	service := NewService()

	if description.Driver == DiscDriver {
		err = service.AddOnDiscIndex(description)
	} else {
		err = service.AddRunTimeIndex(description)
	}

	if err != nil {
		t.Error(err)
	}

	wordsList := []string{"Nissan March", "Honda Fitt", "Wolfsvagen", "Tayota Corolla", "Micra Nissan"}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 5; i++ {
			if description.Driver == DiscDriver {
				err = service.AddOnDiscIndex(description)
			} else {
				err = service.AddRunTimeIndex(description)
			}

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
		for i := 0; i < len(expectedValues); i++ {
			searchConf, _ := NewSearchConfig(wordsList[i], 5, metric.CosineMetric(), 0.7)
			result, err := service.Suggest(description.Name, searchConf)
			if err != nil {
				t.Errorf("Fail suggest %v", err)
			}

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
