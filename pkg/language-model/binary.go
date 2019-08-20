package lm

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alldroll/go-datastructures/rbtree"
	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/mph"
)

// StoreBinaryLMFromGoogleFormat creates a ngram language model from the google ngram format
func StoreBinaryLMFromGoogleFormat(config *Config) error {
	dict, err := buildDictionary(config)

	if err != nil {
		return err
	}

	table, err := buildMPH(dict)

	if err != nil {
		return err
	}

	reader := NewGoogleNGramReader(config.NGramOrder, NewIndexer(dict, table), config.GetOutputPath())
	model, err := reader.Read()

	if err != nil {
		return fmt.Errorf("couldn't read ngrams: %v", err)
	}

	binaryFile, err := os.OpenFile(
		config.GetBinaryPath(),
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return fmt.Errorf("failed to create a binary file: %v", err)
	}

	enc := gob.NewEncoder(binaryFile)

	if err := enc.Encode(&model); err != nil {
		return fmt.Errorf("failed to encode NGramModel in the binary format: %v", err)
	}

	if err := enc.Encode(&table); err != nil {
		return fmt.Errorf("failed to encode MPH in the binary format: %v", err)
	}

	return nil
}

// RetrieveLMFromBinary retrieves a language model from the binary format
func RetrieveLMFromBinary(config *Config) (LanguageModel, error) {
	dict, err := dictionary.OpenCDBDictionary(config.GetDictionaryPath())

	if err != nil {
		return nil, err
	}

	binaryFile, err := os.Open(config.GetBinaryPath())

	if err != nil {
		return nil, fmt.Errorf("Failed to open the lm binary file: %v", err)
	}

	var (
		model NGramModel
		table mph.MPH
		dec   = gob.NewDecoder(binaryFile)
	)

	if err := dec.Decode(&model); err != nil {
		return nil, err
	}

	if err := dec.Decode(&table); err != nil {
		return nil, err
	}

	return NewLanguageModel(model, NewIndexer(dict, table), config), nil
}

// buildDictionary builds a dictionary for the given config
func buildDictionary(config *Config) (dictionary.Dictionary, error) {
	dictReader, err := newDictionaryReader(config)

	if err != nil {
		return nil, err
	}

	dict, err := dictionary.BuildCDBDictionary(dictReader, config.GetDictionaryPath())

	if err != nil {
		return nil, err
	}

	return dict, nil
}

// newDictionaryReader creates an adapter to Iterable interface, that scans all lines
// from the SourcePath and creates pairs of <DocID, Value>
func newDictionaryReader(config *Config) (dictionary.Iterable, error) {
	f, err := os.Open(fmt.Sprintf(fileFormat, config.GetOutputPath(), 1))

	if err != nil {
		return nil, fmt.Errorf("could not open a source file %s", err)
	}

	scanner := bufio.NewScanner(f)

	return &dictionaryReader{
		lineScanner: scanner,
	}, nil
}

// dictionaryReader is an adapter, that implements dictionary.Iterable for bufio.Scanner
type dictionaryReader struct {
	lineScanner *bufio.Scanner
}

type dictItem struct {
	word  dictionary.Value
	count WordCount
}

// Less tells is current elements is bigger than the other
func (n *dictItem) Less(other rbtree.Item) bool {
	o := other.(*dictItem)

	if n.count > o.count {
		return true
	}

	if n.count == o.count {
		return n.word > o.word
	}

	return false
}

// Iterate iterates through each line of the corresponding dictionary
func (dr *dictionaryReader) Iterate(iterator dictionary.Iterator) error {
	tree := rbtree.New()

	for dr.lineScanner.Scan() {
		line := dr.lineScanner.Text()
		tabIndex := strings.Index(line, "\t")
		count, err := strconv.ParseUint(line[tabIndex+1:], 10, 32)

		if err != nil {
			return err
		}

		_, _ = tree.Insert(&dictItem{
			word:  line[:tabIndex],
			count: WordCount(count),
		})
	}

	if err := dr.lineScanner.Err(); err != nil {
		return err
	}

	for docID, iter := dictionary.Key(0), tree.NewIterator(); iter.Next() != nil; docID++ {
		item := iter.Get().(*dictItem)

		if err := iterator(docID, item.word); err != nil {
			return err
		}
	}

	return nil
}

// buildMPH builds a mph from the given dictionary
func buildMPH(dict dictionary.Dictionary) (mph.MPH, error) {
	table, err := mph.BuildMPH(dict)

	if err != nil {
		return nil, fmt.Errorf("failed to build mph: %v", err)
	}

	return table, nil
}
