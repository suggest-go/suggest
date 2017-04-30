package suggest

import (
	"regexp"
	"strings"
)

type cleaner struct {
	reg          *regexp.Regexp
	pad, wrapper string
}

func newCleaner(chars []rune, pad, wrapper string) *cleaner {
	str := string(chars)
	reg, err := regexp.Compile("[^" + str + "]+")
	if err != nil {
		panic(err)
	}

	return &cleaner{reg, pad, wrapper}
}

func (self *cleaner) normalize(word string) string {
	if len(word) < 2 {
		return word
	}

	word = strings.ToLower(word)
	word = strings.Trim(word, " ")
	return word
}

func (self *cleaner) clean(word string) string {
	word = self.normalize(word)
	return self.reg.ReplaceAllString(word, self.pad)
}

func (self *cleaner) wrap(word string) string {
	return self.wrapper + word + self.wrapper
}

func (self *cleaner) cleanAndWrap(word string) string {
	return self.wrap(self.clean(word))
}
