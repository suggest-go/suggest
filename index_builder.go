package suggest

const (
	defaultPad = "$"
	defaultWrap = "$"
	defaultNGramSize = 3
)

// Builder
type Builder interface {
	SetNGramSize(nGramSize int) Builder
	SetAlphabet(alphabet Alphabet) Builder
	SetDictionary(dictionary Dictionary) Builder
	SetPad(pad string) Builder
	SetWrap(wrap string) Builder
	Build() NGramIndex
}

type runTimeBuilderImpl struct {
	nGramSize int
	alphabet Alphabet
	dictionary Dictionary
	pad string
	wrap string
}

func NewRunTimeBuilder() Builder {
	return &runTimeBuilderImpl{
		defaultNGramSize,
		nil,
		nil,
		defaultPad,
		defaultWrap,
	}
}

func (b *runTimeBuilderImpl) SetNGramSize(nGramSize int) Builder {
	b.nGramSize = nGramSize
	return b
}

func (b *runTimeBuilderImpl) SetAlphabet(alphabet Alphabet) Builder {
	b.alphabet = alphabet
	return b
}

func (b *runTimeBuilderImpl) SetDictionary(dictionary Dictionary) Builder {
	b.dictionary = dictionary
	return b
}

func (b *runTimeBuilderImpl) SetPad(pad string) Builder {
	b.pad = pad
	return b
}

func (b *runTimeBuilderImpl) SetWrap(wrap string) Builder {
	b.wrap = wrap
	return b
}

func (b *runTimeBuilderImpl) Build() NGramIndex {
	cleaner := NewCleaner(b.alphabet.Chars(), b.pad, b.wrap)
	generator := NewGenerator(b.nGramSize, b.alphabet)
	indexer := NewIndexer(
		b.nGramSize,
		generator,
		cleaner,
	)

	invertedListsBuilder := NewInMemoryInvertedIndexIndicesBuilder(indexer.Index(b.dictionary))
	invertedIndexIndices := invertedListsBuilder.Build()

	return NewNGramIndex(
		cleaner,
		generator,
		invertedIndexIndices,
	)
}

type builderImpl struct {
	nGramSize int
	alphabet Alphabet
	dictionary Dictionary
	pad string
	wrap string
	pattern string
}

// NewBuilder works with already indexed data
func NewBuilder(pattern string) Builder {
	return &builderImpl{
		defaultNGramSize,
		nil,
		nil,
		defaultPad,
		defaultWrap,
		pattern,
	}
}

func (b *builderImpl) SetNGramSize(nGramSize int) Builder {
	b.nGramSize = nGramSize
	return b
}

func (b *builderImpl) SetAlphabet(alphabet Alphabet) Builder {
	b.alphabet = alphabet
	return b
}

func (b *builderImpl) SetDictionary(dictionary Dictionary) Builder {
	b.dictionary = dictionary
	return b
}

func (b *builderImpl) SetPad(pad string) Builder {
	b.pad = pad
	return b
}

func (b *builderImpl) SetWrap(wrap string) Builder {
	b.wrap = wrap
	return b
}

func (b *builderImpl) Build() NGramIndex {
	cleaner := NewCleaner(b.alphabet.Chars(), b.pad, b.wrap)
	generator := NewGenerator(b.nGramSize, b.alphabet)

	invertedListsBuilder := NewCDBInvertedIndexIndicesBuilder(b.pattern)
	invertedIndexIndices := invertedListsBuilder.Build()

	return NewNGramIndex(
		cleaner,
		generator,
		invertedIndexIndices,
	)
}
