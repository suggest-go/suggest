package suggest

import "sync"

type SuggestService struct {
	sync.RWMutex
	dictionaries map[string]*NGramIndex
	configs      map[string]*Config
}

type Config struct {
	ngramSize int
	topK      int
	name      string
}

func NewConfig(ngramSize int, topK int, name string) *Config {
	return &Config{
		ngramSize,
		topK,
		name,
	}
}

func (self *Config) GetName() string {
	return self.name
}

// Creates an empty SuggestService which uses given metric as "edit distance metric"
func NewSuggestService() *SuggestService {
	return &SuggestService{
		dictionaries: make(map[string]*NGramIndex),
		configs:      make(map[string]*Config),
	}
}

// add/replace new dictionary with given name
func (self *SuggestService) AddDictionary(name string, words []string, config *Config) error {
	ngramIndex := NewNGramIndex(config.ngramSize)
	for _, word := range words {
		ngramIndex.AddWord(word)
	}

	self.Lock()
	self.dictionaries[name] = ngramIndex
	self.configs[name] = config
	self.Unlock()
	return nil
}

// return Top-k approximate strings for given query in dict
func (self *SuggestService) Suggest(dict string, query string) []string {
	self.RLock()
	index, okDict := self.dictionaries[dict]
	config, okConf := self.configs[dict]
	self.RUnlock()
	if okDict && okConf {
		return index.Suggest(query, config.topK)
	}

	return make([]string, 0)
}
