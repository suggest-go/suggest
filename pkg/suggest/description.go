package suggest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/alldroll/suggest/pkg/alphabet"
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

// CreateAlphabet creates a new instance of alphabet from the configuration
func (d *IndexDescription) CreateAlphabet() alphabet.Alphabet {
	return alphabet.CreateAlphabet(d.Alphabet)
}

// GetDictionaryFile returns a path to a dictionary file from the configuration
func (d *IndexDescription) GetDictionaryFile() string {
	return fmt.Sprintf("%s/%s.cdb", d.GetOutputPath(), d.Name)
}

// GetHeaderFile returns a path to a header file from the configuration
func (d *IndexDescription) GetHeaderFile() string {
	return fmt.Sprintf("%s.hd", d.Name)
}

// GetDocumentListFile returns a path to a document list file from the configuration
func (d *IndexDescription) GetDocumentListFile() string {
	return fmt.Sprintf("%s.dl", d.Name)
}

// GetOutputPath returns a output path of the builded index
func (d *IndexDescription) GetOutputPath() string {
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

// CreateWriterConfig creates and returns IndexWriter config from the given index description
func (d *IndexDescription) CreateWriterConfig() index.WriterConfig {
	return index.WriterConfig{
		HeaderFileName:       d.GetHeaderFile(),
		DocumentListFileName: d.GetDocumentListFile(),
	}
}
