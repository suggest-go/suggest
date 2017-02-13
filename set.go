package suggest

type Set struct {
	set  map[string]struct{}
	keys []string
}

func NewSet(values []string) *Set {
	set := make(map[string]struct{}, len(values))
	keys := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := set[value]; !ok {
			set[value] = struct{}{}
			keys = append(keys, value)
		}
	}

	return &Set{set, keys}
}

func (self *Set) GetKeys() []string {
	return self.keys
}

func (self *Set) Contains(key string) bool {
	_, ok := self.set[key]
	return ok
}
