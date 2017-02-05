package suggest

import (
	"bufio"
	"bytes"
	"io"
)

type SuggestService struct {
	dictionaries map[string]*NGramIndex
	topK         int
}

func NewSuggestService(topK int) *SuggestService {
	return &SuggestService{
		make(map[string]*NGramIndex),
		topK,
	}
}

/*
* inspired by https://github.com/jprichardson/readline-go/blob/master/readline.go
 */
func (self *SuggestService) AddDictionary(name string, reader io.Reader) bool {
	if _, ok := self.dictionaries[name]; !ok {
		// TODO throw error
		return false
	}

	ngramIndex := NewNGramIndex(3)
	buf := bufio.NewReader(reader)
	line, err := buf.ReadBytes('\n')
	for err == nil {
		line = bytes.TrimRight(line, "\n")
		if len(line) > 0 {
			if line[len(line)-1] == 13 { //'\r'
				line = bytes.TrimRight(line, "\r")
			}

			ngramIndex.AddWord(string(line))
		}

		line, err = buf.ReadBytes('\n')
	}

	if len(line) > 0 {
		ngramIndex.AddWord(string(line))
	}

	self.dictionaries[name] = ngramIndex
	return true
}

func (self *SuggestService) Suggest(dict string, query string) []string {
	if _, ok := self.dictionaries[dict]; !ok {
		return make([]string, 0)
	}

	return self.dictionaries[dict].Suggest(query, self.topK)
}
