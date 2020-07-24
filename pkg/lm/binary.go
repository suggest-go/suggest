package lm

import (
	"bufio"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/alldroll/go-datastructures/rbtree"
	"github.com/suggest-go/suggest/pkg/dictionary"
	"github.com/suggest-go/suggest/pkg/mph"
	"github.com/suggest-go/suggest/pkg/store"
)

// StoreBinaryLMFromGoogleFormat creates a ngram language model from the google ngram format
func StoreBinaryLMFromGoogleFormat(directory store.Directory, config *Config) error {
	dict, err := buildDictionary(directory, config)

	if err != nil {
		return err
	}

	table, err := buildMPH(dict)

	if err != nil {
		return err
	}

	reader := NewGoogleNGramReader(config.NGramOrder, NewIndexer(dict, table), directory)
	model, err := reader.Read()

	if err != nil {
		return fmt.Errorf("couldn't read ngrams: %w", err)
	}

	out, err := directory.CreateOutput(config.GetBinaryPath())

	if err != nil {
		return fmt.Errorf("failed to create a binary file: %w", err)
	}

	if _, err := model.Store(out); err != nil {
		return fmt.Errorf("failed to encode NGramModel in the binary format: %w", err)
	}

	if _, err := table.Store(out); err != nil {
		return fmt.Errorf("failed to encode MPH in the binary format: %w", err)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("failed to close a binary output: %w", err)
	}

	return nil
}

// RetrieveLMFromBinary retrieves a language model from the binary format
func RetrieveLMFromBinary(directory store.Directory, config *Config) (LanguageModel, error) {
	dict, err := dictionary.OpenCDBDictionary(config.GetDictionaryPath())

	if err != nil {
		return nil, err
	}

	in, err := directory.OpenInput(config.GetBinaryPath())

	if err != nil {
		return nil, fmt.Errorf("failed to open the lm binary file: %w", err)
	}

	var (
		model = &nGramModel{}
		table = mph.New()
	)

	if _, err := model.Load(in); err != nil {
		return nil, err
	}

	if _, err := table.Load(in); err != nil {
		return nil, err
	}

	languageModel, err := NewLanguageModel(model, NewIndexer(dict, table), config)

	// TODO make it clear
	if err == nil {
		runtime.SetFinalizer(languageModel, func(d interface{}) {
			in.Close()
		})
	} else {
		in.Close()
	}

	return languageModel, err
}

// buildDictionary builds a dictionary for the given config
func buildDictionary(directory store.Directory, config *Config) (dictionary.Dictionary, error) {
	dictReader, err := newDictionaryReader(directory)

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
func newDictionaryReader(directory store.Directory) (dictionary.Iterable, error) {
	in, err := directory.OpenInput(fmt.Sprintf(fileFormat, 1))

	if err != nil {
		return nil, fmt.Errorf("could not open a source file %w", err)
	}

	return &dictionaryReader{
		in: in,
	}, nil
}

// dictionaryReader is an adapter, that implements dictionary.Iterable for bufio.Scanner
type dictionaryReader struct {
	in store.Input
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
		return n.word < o.word
	}

	return false
}

// Iterate iterates through each line of the corresponding dictionary
func (dr *dictionaryReader) Iterate(iterator dictionary.Iterator) error {
	tree := rbtree.New()
	lineScanner := bufio.NewScanner(dr.in)

	for lineScanner.Scan() {
		line := lineScanner.Text()
		tabIndex := strings.Index(line, "\t")
		count, err := strconv.ParseUint(line[tabIndex+1:], 10, 32)

		if err != nil {
			return err
		}

		if len(line[:tabIndex]) == 0 {
			continue
		}

		item := &dictItem{
			word:  line[:tabIndex],
			count: WordCount(count),
		}

		_, _ = tree.Insert(item)
	}

	if err := lineScanner.Err(); err != nil {
		return err
	}

	for docID, iter := dictionary.Key(0), tree.NewIterator(); iter.Next() != nil; docID++ {
		item := iter.Get().(*dictItem)

		if err := iterator(docID, item.word); err != nil {
			return fmt.Errorf("failed to iterate through dictionary: %w", err)
		}
	}

	if err := dr.in.Close(); err != nil {
		return fmt.Errorf("failed to close a dictionary input: %w", err)
	}

	return nil
}

// buildMPH builds a mph from the given dictionary
func buildMPH(dict dictionary.Dictionary) (mph.MPH, error) {
	table := mph.New()

	if err := table.Build(dict); err != nil {
		return nil, fmt.Errorf("failed to build a mph table: %w", err)
	}

	return table, nil
}
