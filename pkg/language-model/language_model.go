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

// languageModel implements LanguageModel interface
type languageModel struct {
	model   NGramModel
	indexer Indexer
	config  *Config
}

// NewLanguageModel creates a new instance of a LanguageModel
func NewLanguageModel(
	model NGramModel,
	indexer Indexer,
	config *Config,
) LanguageModel {
	return &languageModel{
		model:   model,
		indexer: indexer,
		config:  config,
	}
}

// ScoreSentence scores and returns a weight in the language model for the given sentence
func (lm *languageModel) ScoreSentence(sentence Sentence) (float64, error) {
	ids, err := MapIntoListOfWordIDs(lm, lm.wrapSentence(sentence))

	if err != nil {
		return 0, err
	}

	return lm.ScoreWordIDs(ids), nil
}

// ScoreWordIDs scores and returns a lm weight for the given sequence of WordID
func (lm *languageModel) ScoreWordIDs(sequence []WordID) float64 {
	score := 0.0

	for _, nGrams := range lm.split(sequence) {
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
func (lm *languageModel) wrapSentence(sentence Sentence) Sentence {
	return lm.leftWrapSentence(lm.rightWrapSentence(sentence))
}

// leftWrapSentence prepends the start symbol to the given sentence
func (lm *languageModel) leftWrapSentence(sentence Sentence) Sentence {
	return append([]Token{lm.config.StartSymbol}, sentence...)
}

// rightWrapSentence appends the end symbol to the given sentence
func (lm *languageModel) rightWrapSentence(sentence Sentence) Sentence {
	return append(sentence, lm.config.EndSymbol)
}
