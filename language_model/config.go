package language_model

import "io"

type Config struct {
	Name       string   `json:"name"`
	NGramOrder uint8    `json:"nGramOrder"`
	SourcePath string   `json:"source"`
	OutputPath string   `json:"output"`
	Alphabet   []string `json:"alphabet"`
	Separators []string `json:"separators"`
}

func ReadConfig(reader io.Reader) ([]Config, error) {
	return nil, nil
}
