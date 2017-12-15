package main

import (
	"io"
	"encoding/json"
	"io/ioutil"
)

type IndexConfig struct {
	Name string `json:"name"`
	NGramSize int `json:"nGramSize"`
	SourcePath string `json:"source"`
	OutputPath string `json:"output"`
	Alphabet []string `json:"alphabet"`
	Pad string `json:"pad"`
	//Wrap [2]string `json:"wrap"`
	Wrap string `json:"wrap"`
}

func readConfig(reader io.Reader) ([]IndexConfig, error) {
	var configs []IndexConfig

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
