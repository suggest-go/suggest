package lm

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/alldroll/suggest/pkg/alphabet"
)

// Config represents a configuration of a language model
type Config struct {
	Name        string   `json:"name"`
	NGramOrder  uint8    `json:"nGramOrder"`
	SourcePath  string   `json:"source"`
	OutputPath  string   `json:"output"`
	Alphabet    []string `json:"alphabet"`
	Separators  []string `json:"separators"`
	StartSymbol string   `json:"startSymbol"`
	EndSymbol   string   `json:"endSymbol"`
}

// GetWordsAlphabet returns a word alphabet corresponding to the declaration
func (c *Config) GetWordsAlphabet() alphabet.Alphabet {
	return alphabet.CreateAlphabet(c.Alphabet)
}

// GetSeparatorsAlphabet returns a separators alphabet corresponding to the declaration
func (c *Config) GetSeparatorsAlphabet() alphabet.Alphabet {
	return alphabet.CreateAlphabet(c.Separators)
}

// GetDictionaryPath returns a stored path for the dictionary
func (c *Config) GetDictionaryPath() string {
	return fmt.Sprintf("%s/%s.cdb", c.OutputPath, c.Name)
}

// GetMPHPath returns a stored path for the mph
func (c *Config) GetMPHPath() string {
	return fmt.Sprintf("%s/%s.mph", c.OutputPath, c.Name)
}

// GetBinaryPath returns a stored path for the binary lm
func (c *Config) GetBinaryPath() string {
	return fmt.Sprintf("%s/%s.lm", c.OutputPath, c.Name)
}

// ReadConfig reads a language model config from the given reader
func ReadConfig(reader io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
