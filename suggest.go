package suggest

type SuggestService struct {
	dictionaries map[string]*NGramIndex
	ngramSize    int
	editDistance EditDistance
}

func NewSuggestService(ngramSize, metric int) *SuggestService {
	editDistance, err := GetEditDistance(metric, ngramSize)
	if err != nil {
		panic(err)
	}

	return &SuggestService{
		make(map[string]*NGramIndex),
		ngramSize,
		editDistance,
	}
}

func (self *SuggestService) AddDictionary(name string, words []string) bool {
	if _, ok := self.dictionaries[name]; ok {
		//"dictionary already exists" /* TODO log me */
		return false
	}

	ngramIndex := NewNGramIndex(self.ngramSize, self.editDistance)
	for _, word := range words {
		ngramIndex.AddWord(word)
	}

	self.dictionaries[name] = ngramIndex
	return true
}

func (self *SuggestService) Suggest(dict string, query string, topK int) []string {
	if _, ok := self.dictionaries[dict]; !ok {
		return make([]string, 0)
	}

	return self.dictionaries[dict].Suggest(query, topK)
}
