package suggest

import "github.com/alldroll/rbtree"

type set struct {
	tree *rbtree.RBTree
}

func newSet() *set {
	return &set{rbtree.New()}
}

type ngramProfile struct {
	ngram     string
	frequency int
}

func (self *ngramProfile) Less(other rbtree.Item) bool {
	return self.ngram < other.(*ngramProfile).ngram
}

func (self *set) Add(p *ngramProfile) {
	self.tree.Insert(p)
}

func (self *set) Size() int {
	return self.tree.Len()
}

// TODO fix monkeycode
func union(a, b *set) []*ngramProfile {
	aIter := a.tree.NewIterator()
	bIter := b.tree.NewIterator()
	aIter.Next()
	bIter.Next()
	result := make([]*ngramProfile, 0, a.Size()+b.Size())
	for {
		if !aIter.IsValid() {
			for bIter.IsValid() {
				result = append(result, bIter.Get().(*ngramProfile))
				bIter.Next()
			}

			break
		}

		if !bIter.IsValid() {
			for aIter.IsValid() {
				result = append(result, aIter.Get().(*ngramProfile))
				aIter.Next()
			}

			break
		}

		if aIter.Get().Less(bIter.Get()) {
			result = append(result, aIter.Get().(*ngramProfile))
			aIter.Next()
		} else if bIter.Get().Less(aIter.Get()) {
			result = append(result, bIter.Get().(*ngramProfile))
			bIter.Next()
		} else {
			result = append(result, bIter.Get().(*ngramProfile))
			bIter.Next()
			aIter.Next()
		}
	}

	return result
}

// TODO fix monkeycode
func intersection(a, b *set) []*ngramProfile {
	aIter := a.tree.NewIterator()
	bIter := b.tree.NewIterator()
	aIter.Next()
	bIter.Next()
	result := make([]*ngramProfile, 0, b.Size())
	for aIter.IsValid() && bIter.IsValid() {
		if aIter.Get().Less(bIter.Get()) {
			aIter.Next()
		} else if bIter.Get().Less(aIter.Get()) {
			bIter.Next()
		} else {
			result = append(result, aIter.Get().(*ngramProfile))
			bIter.Next()
			aIter.Next()
		}
	}

	return result
}
