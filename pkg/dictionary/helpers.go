package dictionary

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/alldroll/cdb"
)

// BuildCDBDictionary is a helper for building a CDB dictionary from the sourcePath
// Saves the dictionary to destinationPath
func BuildCDBDictionary(sourcePath, destinationPath string) (Dictionary, error) {
	sourceFile, err := os.OpenFile(sourcePath, os.O_RDONLY, 0)

	if err != nil {
		return nil, fmt.Errorf("Unable to open source file %s", err)
	}

	destinationFile, err := os.OpenFile(
		destinationPath,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		0644,
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create dictionary file %s", err)
	}

	cdbHandle := cdb.New()
	cdbWriter, err := cdbHandle.GetWriter(destinationFile)

	if err != nil {
		return nil, fmt.Errorf("Failed to create cdb writer %s", err)
	}

	var (
		docID   Key
		key     = make([]byte, 4)
		scanner = bufio.NewScanner(sourceFile)
	)

	for scanner.Scan() {
		binary.LittleEndian.PutUint32(key, docID)

		if err := cdbWriter.Put(key, scanner.Bytes()); err != nil {
			return nil, fmt.Errorf("Failed to put record to cdb %s", err)
		}

		docID++
	}

	if err := cdbWriter.Close(); err != nil {
		return nil, fmt.Errorf("Failed to save cdb dictionary %s", err)
	}

	return NewCDBDictionary(destinationFile)
}
