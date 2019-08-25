package suggest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/alldroll/suggest/pkg/analysis"
	"github.com/alldroll/suggest/pkg/index"
)

// Driver represents storage type of an inverted index
type Driver string

const (
	// RAMDriver means that an inverted index is stored in RAM
	RAMDriver Driver = "RAM"
	// DiscDriver means that an inverted index is stored on FS and was indexed before
	DiscDriver Driver = "DISC"
)

// IndexDescription is config for NgramIndex structure
type IndexDescription struct {
	Driver     Driver    `json:"driver"`
	Name       string    `json:"name"`
	NGramSize  int       `json:"nGramSize"`
	SourcePath string    `json:"source"`
	OutputPath string    `json:"output"`
	Alphabet   []string  `json:"alphabet"`
	Pad        string    `json:"pad"`
	Wrap       [2]string `json:"wrap"`
	basePath   string
}

// GetDictionaryFile returns a path to a dictionary file from the configuration
func (d *IndexDescription) GetDictionaryFile() string {
	return fmt.Sprintf("%s/%s.cdb", d.GetIndexPath(), d.Name)
}

// GetIndexPath returns a output path of the built index
func (d *IndexDescription) GetIndexPath() string {
	if !path.IsAbs(d.OutputPath) {
		return fmt.Sprintf("%s/%s", d.basePath, d.OutputPath)
	}

	return d.OutputPath
}

// GetSourcePath returns a source path of the index description
func (d *IndexDescription) GetSourcePath() string {
	if !path.IsAbs(d.SourcePath) {
		return fmt.Sprintf("%s/%s", d.basePath, d.SourcePath)
	}

	return d.SourcePath
}

// GetWriterConfig creates and returns IndexWriter config from the given index description
func (d *IndexDescription) GetWriterConfig() index.WriterConfig {
	return index.WriterConfig{
		HeaderFileName:       d.getHeaderFile(),
		DocumentListFileName: d.getDocumentListFile(),
	}
}

// GetIndexTokenizer returns a tokenizer for indexing
func (d *IndexDescription) GetIndexTokenizer() analysis.Tokenizer {
	return NewSuggestTokenizer(*d)
}

// getHeaderFile returns a path to a header file from the configuration
func (d *IndexDescription) getHeaderFile() string {
	return fmt.Sprintf("%s.hd", d.Name)
}

// getDocumentListFile returns a path to a document list file from the configuration
func (d *IndexDescription) getDocumentListFile() string {
	return fmt.Sprintf("%s.dl", d.Name)
}

// ReadConfigs reads and returns a list of IndexDescription from the given reader
func ReadConfigs(configPath string) ([]IndexDescription, error) {
	configFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer configFile.Close()
	var configs []IndexDescription

	data, err := ioutil.ReadAll(configFile)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, err
	}

	basePath := path.Dir(configPath)

	for i, c := range configs {
		c.basePath = basePath
		configs[i] = c
	}

	return configs, nil
}
