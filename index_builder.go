package suggest

const (
	defaultPad = "$"
	defaultWrap = "$"
	defaultNGramSize = 3
)

type Builder interface {
	SetNGramSize(nGramSize int) Builder
	SetAlphabet(alphabet Alphabet) Builder
	SetDictionary(dictionary Dictionary) Builder
	SetPad(pad string) Builder
	SetWrap(wrap string) Builder
	Build() NGramIndex
}

type builderImpl struct {
	nGramSize int
	alphabet Alphabet
	dictionary Dictionary
	pad string
	wrap string
}

func NewBuilder() Builder {
	return &builderImpl{
		defaultNGramSize,
		nil,
		nil,
		defaultPad,
		defaultWrap,
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

	invertedListsBuilder := NewInvertedListBuilder(
		b.nGramSize,
		b.dictionary,
		generator,
		cleaner,
	)

	invertedListIndices := invertedListsBuilder.Build()

	return NewNGramIndex(
		cleaner,
		generator,
		invertedListIndices,
	)
}
