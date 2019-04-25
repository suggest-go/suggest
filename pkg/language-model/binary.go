package lm

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strings"

	"github.com/alldroll/suggest/pkg/dictionary"
	"github.com/alldroll/suggest/pkg/mph"
)

// StoreBinaryLMFromGoogleFormat creates a ngram language model from the google ngram format
func StoreBinaryLMFromGoogleFormat(config *Config) error {
	dict, err := buildDictionary(config)

	if err != nil {
		return err
	}

	mph, err := buildMPH(dict)

	if err != nil {
		return err
	}

	reader := NewGoogleNGramReader(config.NGramOrder, NewIndexer(dict, mph), config.OutputPath)
	model, err := reader.Read()

	if err != nil {
		return fmt.Errorf("Couldn't read ngrams: %v", err)
	}

	binaryFile, err := os.OpenFile(
		config.GetBinaryPath(),
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return fmt.Errorf("Failed to create a binary file: %v", err)
	}

	enc := gob.NewEncoder(binaryFile)

	if err := enc.Encode(&model); err != nil {
		return fmt.Errorf("Failed to encode NGramModel in the binary format: %v", err)
	}

	if err := enc.Encode(&mph); err != nil {
		return fmt.Errorf("Failed to encode MPH in the binary format: %v", err)
	}

	return nil
}

// RetrieveLMFromBinary retrives a language model from the binary format
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
		mph   mph.MPH
		dec   = gob.NewDecoder(binaryFile)
	)

	if err := dec.Decode(&model); err != nil {
		return nil, err
	}

	if err := dec.Decode(&mph); err != nil {
		return nil, err
	}

	return NewLanguageModel(model, NewIndexer(dict, mph), config), nil
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
	f, err := os.Open(fmt.Sprintf(fileFormat, config.OutputPath, 1))

	if err != nil {
		return nil, fmt.Errorf("Could not open a source file %s", err)
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

// Iterate iterates through each line of the corresponding dictionary
func (dr *dictionaryReader) Iterate(iterator dictionary.Iterator) error {
	docID := dictionary.Key(0)

	for dr.lineScanner.Scan() {
		line := dr.lineScanner.Text()
		tabIndex := strings.Index(line, "\t")

		if err := iterator(docID, line[:tabIndex]); err != nil {
			return err
		}

		docID++
	}

	return dr.lineScanner.Err()
}

// buildMPH builds a mph from the given dictionary
func buildMPH(dict dictionary.Dictionary) (mph.MPH, error) {
	mph, err := mph.BuildMPH(dict)

	if err != nil {
		return nil, fmt.Errorf("Failed to build mph: %v", err)
	}

	return mph, nil
}
