package suggest

import (
	"encoding/json"
	"fmt"
	"github.com/alldroll/suggest/alphabet"
	"io"
	"io/ioutil"
)

// IndexDescription is config for NgramIndex structure
type IndexDescription struct {
	Name       string    `json:"name"`
	NGramSize  int       `json:"nGramSize"`
	SourcePath string    `json:"source"`
	OutputPath string    `json:"output"`
	Alphabet   []string  `json:"alphabet"`
	Pad        string    `json:"pad"`
	Wrap       [2]string `json:"wrap"`
}

func ReadConfigs(reader io.Reader) ([]IndexDescription, error) {
	var configs []IndexDescription

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (d *IndexDescription) CreateAlphabet() alphabet.Alphabet {
	return alphabet.CreateAlphabet(d.Alphabet)
}

func (d *IndexDescription) GetDictionaryFile() string {
	return fmt.Sprintf("%s/%s.cdb", d.OutputPath, d.Name)
}

func (d *IndexDescription) GetHeaderFile() string {
	return fmt.Sprintf("%s/%s.hd", d.OutputPath, d.Name)
}

func (d *IndexDescription) GetDocumentListFile() string {
	return fmt.Sprintf("%s/%s.dl", d.OutputPath, d.Name)
}
