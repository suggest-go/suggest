package suggest

import "github.com/alldroll/rbtree"

type WordProfile struct {
	word   string
	ngrams []NGram
}

func (self *WordProfile) GetWord() string {
	return self.word
}

// Return unique ordered NGram[]
func (self *WordProfile) GetNGrams() []NGram {
	return self.ngrams
}

type NGram interface {
	GetValue() string
	GetFrequency() int
}

func GetWordProfile(word string, k int) *WordProfile {
	ngrams := SplitIntoNGrams(word, k)
	tree := rbtree.New()
	for _, ngram := range ngrams {
		n := &ngramImp{ngram, 0}
		item := tree.Get(n)
		if item == nil {
			tree.Insert(n)
		} else {
			n = item.(*ngramImp)
		}

		n.frequency++
	}

	list := make([]NGram, 0, tree.Len())
	iter := tree.NewIterator()
	for {
		item := iter.Next()
		if item == nil {
			break
		}

		list = append(list, item.(*ngramImp))
	}

	return &WordProfile{word, list}
}

type ngramImp struct {
	value     string
	frequency int
}

func (self *ngramImp) GetValue() string {
	return self.value
}

func (self *ngramImp) GetFrequency() int {
	return self.frequency
}

func (self *ngramImp) Less(other rbtree.Item) bool {
	return self.value < other.(*ngramImp).value
}
