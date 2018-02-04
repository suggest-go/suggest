package suggest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

var (
	alphabetMap = map[string]Alphabet{
		"english": NewEnglishAlphabet(),
		"russian": NewRussianAlphabet(),
		"numbers": NewNumberAlphabet(),
	}
)

// IndexDescription is config for NgramIndex structure
type IndexDescription struct {
	Name       string   `json:"name"`
	NGramSize  int      `json:"nGramSize"`
	SourcePath string   `json:"source"`
	OutputPath string   `json:"output"`
	Alphabet   []string `json:"alphabet"`
	Pad        string   `json:"pad"`
	Wrap       string   `json:"wrap"`
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

func (d *IndexDescription) CreateAlphabet() Alphabet {
	alphabets := make([]Alphabet, 0)

	for _, symbols := range d.Alphabet {
		if alphabet, ok := alphabetMap[symbols]; ok {
			alphabets = append(alphabets, alphabet)
			continue
		}

		alphabets = append(alphabets, NewSimpleAlphabet([]rune(symbols)))
	}

	return NewCompositeAlphabet(alphabets)
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
