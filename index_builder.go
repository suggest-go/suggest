package suggest

const (
	defaultPad       = "$"
	defaultWrap      = "$"
	defaultNGramSize = 3
)

// Builder
type Builder interface {
	// SetNGramSize
	SetNGramSize(nGramSize int) Builder
	// SetAlphabet
	SetAlphabet(alphabet Alphabet) Builder
	// SetDictionary
	SetDictionary(dictionary Dictionary) Builder
	// SetPad
	SetPad(pad string) Builder
	// SetWrap
	SetWrap(wrap string) Builder
	// Build
	Build() NGramIndex
}

// runTimeBuilderImpl implements Builder interface
type runTimeBuilderImpl struct {
	nGramSize  int
	alphabet   Alphabet
	dictionary Dictionary
	pad        string
	wrap       string
}

// NewRunTimeBuilder returns new instance of runTimeBuilderImpl
func NewRunTimeBuilder() Builder {
	return &runTimeBuilderImpl{
		nGramSize:  defaultNGramSize,
		alphabet:   nil,
		dictionary: nil,
		pad:        defaultPad,
		wrap:       defaultWrap,
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
		&CPMerge{},
	)
}

// builderImpl implements Builder interface
type builderImpl struct {
	nGramSize        int
	alphabet         Alphabet
	dictionary       Dictionary
	pad              string
	wrap             string
	headerPath       string
	documentListPath string
}

// NewBuilder works with already indexed data
func NewBuilder(headerPath, documentListPath string) Builder {
	return &builderImpl{
		nGramSize:        defaultNGramSize,
		alphabet:         nil,
		dictionary:       nil,
		pad:              defaultPad,
		wrap:             defaultWrap,
		headerPath:       headerPath,
		documentListPath: documentListPath,
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

	invertedListsBuilder := NewOnDiscInvertedIndexIndicesBuilder(b.headerPath, b.documentListPath)
	invertedIndexIndices := invertedListsBuilder.Build()

	return NewNGramIndex(
		cleaner,
		generator,
		invertedIndexIndices,
		&CPMerge{},
	)
}
