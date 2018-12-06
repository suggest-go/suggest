package vgram

import (
	"github.com/alldroll/suggest/pkg/dictionary"
	"reflect"
	"testing"
)

func TestBuildFrequencyTree(t *testing.T) {
	dictionary := dictionary.NewInMemoryDictionary([]string{
		"stick",
		"stich",
		"such",
		"stuck",
	})

	builder := &VGramDictionaryBuilder{
		qMin:       2,
		qMax:       4,
		threshold:  2,
		dictionary: dictionary,
	}

	type data struct {
		nodeFrequency   uint32
		markerExists    bool
		markerFrequency uint32
	}

	expected := map[string]data{
		"c":    {4, false, 0},
		"i":    {2, false, 0},
		"s":    {4, false, 0},
		"t":    {3, false, 0},
		"u":    {2, false, 0},
		"ch":   {2, true, 2},
		"ck":   {2, true, 2},
		"ic":   {2, true, 0},
		"st":   {3, true, 0},
		"su":   {1, true, 0},
		"ti":   {2, true, 0},
		"tu":   {1, true, 0},
		"uc":   {2, true, 0},
		"suc":  {1, true, 0},
		"sti":  {2, true, 0},
		"stu":  {1, true, 0},
		"tic":  {2, true, 0},
		"tuc":  {1, true, 0},
		"ick":  {1, true, 1},
		"ich":  {1, true, 1},
		"uck":  {1, true, 1},
		"uch":  {1, true, 1},
		"stic": {2, true, 2},
		"stuc": {1, true, 1},
		"such": {1, true, 1},
		"tich": {1, true, 1},
		"tick": {1, true, 1},
		"tuck": {1, true, 1},
	}

	actual := make(map[string]data, 0)
	trie := builder.buildFrequencyTrie()

	trie.Walk(func(key string, node Node) {
		markerFreq := uint32(0)
		marker := node.GetMarker()

		if marker != nil {
			markerFreq = marker.GetFrequency()
		}

		actual[key] = data{
			nodeFrequency:   node.GetFrequency(),
			markerExists:    marker != nil,
			markerFrequency: markerFreq,
		}

		if !reflect.DeepEqual(actual[key], expected[key]) {
			t.Errorf("%s expected %v, got %v", key, expected[key], actual[key])
		}
	})

	if len(actual) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
