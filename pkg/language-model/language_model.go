package lm

// LanguageModel is an interface for an n-gram language model
type LanguageModel interface {
	// ScoreSentence scores and returns a lm weight for the given sentence
	ScoreSentence(sentence Sentence) (float64, error)
	// ScoreWordIDs scores and returns a lm weight for the given sequence of nGrams
	ScoreWordIDs(sequence []WordID) float64
	// GetWordID returns id for the given token
	GetWordID(token Token) (WordID, error)
}

// languageModel implements LanguageModel interface
type languageModel struct {
	model       NGramModel
	indexer     Indexer
	config      *Config
	startSymbol WordID
	endSymbol   WordID
}

// NewLanguageModel creates a new instance of a LanguageModel
func NewLanguageModel(
	model NGramModel,
	indexer Indexer,
	config *Config,
) LanguageModel {
	// TODO handle me
	startSymbol, _ := indexer.Get(config.StartSymbol)
	endSymbol, _ := indexer.Get(config.EndSymbol)

	return &languageModel{
		model:       model,
		indexer:     indexer,
		config:      config,
		startSymbol: startSymbol,
		endSymbol:   endSymbol,
	}
}

// MapIntoListOfWordIDs maps the given sentence into a list of WordIDs
func MapIntoListOfWordIDs(lm LanguageModel, sentence Sentence) ([]WordID, error) {
	ids := make([]WordID, 0, len(sentence))

	for _, token := range sentence {
		index, err := lm.GetWordID(token)

		if err != nil {
			return nil, err
		}

		ids = append(ids, index)
	}

	return ids, nil
}

// ScoreSentence scores and returns a weight in the language model for the given sentence
func (lm *languageModel) ScoreSentence(sentence Sentence) (float64, error) {
	ids, err := MapIntoListOfWordIDs(lm, sentence)

	if err != nil {
		return 0, err
	}

	return lm.ScoreWordIDs(ids), nil
}

// ScoreWordIDs scores and returns a lm weight for the given sequence of WordID
func (lm *languageModel) ScoreWordIDs(sequence []WordID) float64 {
	score := 0.0

	for _, nGrams := range lm.split(lm.wrapSentence(sequence)) {
		score += lm.model.Score(nGrams)
	}

	return score
}

// GetWordID returns id for the given token
func (lm *languageModel) GetWordID(token Token) (WordID, error) {
	return lm.indexer.Get(token)
}

// split splits the given sequence of WordIDs to nGrams
func (lm *languageModel) split(sequence []WordID) NGrams {
	return SplitIntoNGrams(sequence, lm.config.NGramOrder)
}

// wrapSentence wraps the given sentence with start and end symbols
func (lm *languageModel) wrapSentence(sentence []WordID) []WordID {
	return lm.leftWrapSentence(lm.rightWrapSentence(sentence))
}

// leftWrapSentence prepends the start symbol to the given sentence
func (lm *languageModel) leftWrapSentence(sentence []WordID) []WordID {
	return append([]WordID{lm.startSymbol}, sentence...)
}

// rightWrapSentence appends the end symbol to the given sentence
func (lm *languageModel) rightWrapSentence(sentence []WordID) []WordID {
	return append(sentence, lm.endSymbol)
}
