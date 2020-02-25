package dictionary

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/alldroll/cdb"
	"github.com/suggest-go/suggest/pkg/utils"
)

// OpenCDBDictionary opens a dictionary from cdb file
func OpenCDBDictionary(path string) (Dictionary, error) {
	dictionaryFile, err := utils.NewMMapReader(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open cdb dictionary file: %w", err)
	}

	return NewCDBDictionary(dictionaryFile)
}

// OpenRAMDictionary opens a dictionary from the given path and stores items in RAM
func OpenRAMDictionary(path string) (dict Dictionary, err error) {
	dictionaryFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open dictionary file: %w", err)
	}

	defer func() {
		if cErr := dictionaryFile.Close(); cErr != nil {
			err = cErr
		}
	}()

	scanner := bufio.NewScanner(dictionaryFile)
	collection := make([]string, 0)

	for scanner.Scan() {
		collection = append(collection, scanner.Text())
	}

	dict = NewInMemoryDictionary(collection)

	return
}

// BuildCDBDictionary is a helper for building a CDB dictionary from the sourcePath
// Saves the dictionary to destinationPath
func BuildCDBDictionary(iterator Iterable, destinationPath string) (Dictionary, error) {
	destinationFile, err := os.OpenFile(
		destinationPath,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create dictionary file %w", err)
	}

	cdbHandle := cdb.New()
	cdbWriter, err := cdbHandle.GetWriter(destinationFile)

	if err != nil {
		return nil, fmt.Errorf("failed to create cdb writer %w", err)
	}

	key := make([]byte, 4)

	err = iterator.Iterate(func(docID Key, word Value) error {
		binary.LittleEndian.PutUint32(key, docID)

		if err := cdbWriter.Put(key, []byte(word)); err != nil {
			return fmt.Errorf("failed to put record to cdb: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate through a dictionary: %w", err)
	}

	if err := cdbWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to save cdb dictionary %w", err)
	}

	if err := destinationFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close cdb file %w", err)
	}

	return OpenCDBDictionary(destinationPath)
}
