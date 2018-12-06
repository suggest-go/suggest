package language_model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Config struct {
	NGramOrder uint8    `json:"nGramOrder"`
	SourcePath string   `json:"source"`
	OutputPath string   `json:"output"`
	Alphabet   []string `json:"alphabet"`
	Separators []string `json:"separators"`
}

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
