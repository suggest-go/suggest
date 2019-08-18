package lm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

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
	basePath    string
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
	return fmt.Sprintf("%s/%s.cdb", c.GetOutputPath(), c.Name)
}

// GetMPHPath returns a stored path for the mph
func (c *Config) GetMPHPath() string {
	return fmt.Sprintf("%s/%s.mph", c.GetOutputPath(), c.Name)
}

// GetBinaryPath returns a stored path for the binary lm
func (c *Config) GetBinaryPath() string {
	return fmt.Sprintf("%s/%s.lm", c.GetOutputPath(), c.Name)
}

// GetOutputPath returns a output path of the builded index
func (c *Config) GetOutputPath() string {
	if !path.IsAbs(c.OutputPath) {
		return fmt.Sprintf("%s/%s", c.basePath, c.OutputPath)
	}

	return c.OutputPath
}

// GetSourcePath returns a source path of the index description
func (c *Config) GetSourcePath() string {
	if !path.IsAbs(c.SourcePath) {
		return fmt.Sprintf("%s/%s", c.basePath, c.SourcePath)
	}

	return c.SourcePath
}

// ReadConfig reads a language model config from the given reader
func ReadConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)

	if err != nil {
		return nil, err
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	config.basePath = path.Dir(configPath)

	return &config, nil
}
