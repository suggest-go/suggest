package lm

// LanguageModel is an interface for an n-gram language model
type LanguageModel interface {
	// ScoreSentence scores and returns a weight in the language model for the given sentence
	ScoreSentence(sentence Sentence) float64
	// Predict returns possible next tokens for the given sentence
	Predict(sentence Sentence, topK int) []Token
}

// languageModel implements LanguageModel interface
type languageModel struct {
	model     NGramModel
	generator Generator
	indexer   Indexer
}

// NewLanguageModel creates LanguageModel instance
func NewLanguageModel(model NGramModel, generator Generator, indexer Indexer) LanguageModel {
	return &languageModel{
		model:     model,
		generator: generator,
		indexer:   indexer,
	}
}

// ScoreSentence scores and returns a weight in the language model for the given sentence
func (lm *languageModel) ScoreSentence(sentence Sentence) float64 {
	score := 0.0
	nGramsIds := make([]WordID, 0, 8)

	for _, nGrams := range lm.generator.Generate(sentence) {
		for _, nGram := range nGrams {
			nGramsIds = append(nGramsIds, lm.indexer.Get(nGram))
		}

		score += lm.model.Score(nGramsIds)
		nGramsIds = nGramsIds[:0]
	}

	return score
}

// Predict returns possible next tokens for the given sentence
func (lm *languageModel) Predict(sentence Sentence, topK int) []Token {
	nGrams := lm.generator.Generate(sentence)
	if len(nGrams) == 0 {
		return []Token{}
	}

	last := nGrams[len(nGrams)-1]
	nGramsIds := make([]WordID, 0, 8)

	for _, nGram := range last {
		nGramsIds = append(nGramsIds, lm.indexer.Get(nGram))
	}

	nextWordIDs, err := lm.model.Next(nGramsIds[len(nGramsIds)-1:])
	if err != nil {
		panic(err)
	}

	if len(nextWordIDs) == 0 {
		return []Token{}
	}

	// TODO implement me!

	return nil
}
