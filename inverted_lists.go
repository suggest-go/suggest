package suggest

type invertedListsT map[int][]int

type InvertedListsIndices interface {
	IndexToWordKey(id int) WordKey
	Get(index int) invertedListsT
	Size() int
}

type InvertedListsBuilder interface {
	Build() InvertedListsIndices
}

type invertedListIndicesImpl struct {
	indices []invertedListsT
	tokens []WordKey
}

func (i *invertedListIndicesImpl) IndexToWordKey(id int) WordKey {
	if id >= 0 && int(id) < len(i.tokens) {
		return i.tokens[id]
	}

	return nil
}

func (i *invertedListIndicesImpl) Get(index int) invertedListsT {
	if index >= 0 && index < len(i.indices) {
		return i.indices[index]
	}

	return nil
}

func (i *invertedListIndicesImpl) Size() int {
	return len(i.indices)
}

func NewInvertedListBuilder(
	nGramSize int,
	dictionary Dictionary,
	generator Generator,
	cleaner Cleaner,
) InvertedListsBuilder {
	return &invertedListBuilderImpl{
		nGramSize,
		dictionary,
		generator,
		cleaner,
		make([]invertedListsT, 0),
		make([]WordKey, 0),
	}
}

type invertedListBuilderImpl struct {
	nGramSize int
	dictionary Dictionary
	generator Generator
	cleaner Cleaner
	indices []invertedListsT
	tokens []WordKey
}

func (b *invertedListBuilderImpl) Build() InvertedListsIndices {
	i := b.dictionary.Iterator()
	for {
		key, word := i.GetPair()
		if len(word) >= b.nGramSize {
			b.addWord(word, key)
		}

		if !i.Next() {
			break
		}
	}

	return &invertedListIndicesImpl{
		b.indices,
		b.tokens,
	}
}

func (b *invertedListBuilderImpl) addWord(word string, key WordKey) {
	prepared := b.cleaner.Clean(word)
	set := b.generator.Generate(prepared)
	cardinality := len(set)

	if len(b.indices) <= cardinality {
		tmp := make([]invertedListsT, cardinality+1, cardinality*2)
		copy(tmp, b.indices)
		b.indices = tmp
	}

	invertedLists := b.indices[cardinality]
	if invertedLists == nil {
		invertedLists = make(invertedListsT)
		b.indices[cardinality] = invertedLists
	}

	keyToIndex := len(b.tokens)
	for _, index := range set {
		invertedLists[index] = append(invertedLists[index], keyToIndex)
	}

	b.tokens = append(b.tokens, key)
}
