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

func (c *cleaner) normalize(word string) string {
	if len(word) < 2 {
		return word
	}

	word = strings.ToLower(word)
	word = strings.Trim(word, " ")
	return word
}

func (c *cleaner) clean(word string) string {
	word = c.normalize(word)
	return c.reg.ReplaceAllString(word, c.pad)
}

func (c *cleaner) wrap(word string) string {
	return c.wrapper + word + c.wrapper
}

func (c *cleaner) cleanAndWrap(word string) string {
	return c.wrap(c.clean(word))
}
