package suggest

import (
	"regexp"
	"strings"
)

// Cleaner provides api for preparing word for index search
type Cleaner interface {
	// Clean returns prepared "cleaned" string
	Clean(word string) string
}

type cleanerImpl struct {
	reg          *regexp.Regexp
	pad, wrapper string
}

// NewCleaner returns new Cleaner object
func NewCleaner(chars []rune, pad, wrapper string) Cleaner {
	str := string(chars)
	reg, err := regexp.Compile("[^" + str + "]+")
	if err != nil {
		panic(err)
	}

	return &cleanerImpl{reg, pad, wrapper}
}

func (c *cleanerImpl) Clean(word string) string {
	return c.wrap(c.clean(word))
}

func (c *cleanerImpl) normalize(word string) string {
	if len(word) < 2 {
		return word
	}

	word = strings.ToLower(word)
	word = strings.Trim(word, " ")
	return word
}

func (c *cleanerImpl) clean(word string) string {
	word = c.normalize(word)
	return c.reg.ReplaceAllString(word, c.pad)
}

func (c *cleanerImpl) wrap(word string) string {
	return c.wrapper + word + c.wrapper
}
