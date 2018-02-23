package index

import (
	"regexp"
	"strings"
)

// Cleaner provides api for preparing word for index search
type Cleaner interface {
	// Clean returns prepared "cleaned" string
	Clean(word string) string
}

// cleanerImpl implements Cleaner interface
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

	return &cleanerImpl{
		reg:     reg,
		pad:     pad,
		wrapper: wrapper,
	}
}

// Clean returns prepared "cleaned" string
func (c *cleanerImpl) Clean(word string) string {
	return c.wrap(c.clean(word))
}

// normalize returns normalized word:
// * to lower case
// * trim spaces
func (c *cleanerImpl) normalize(word string) string {
	word = strings.ToLower(word)
	word = strings.Trim(word, " ")
	return word
}

// clean normalize word and replace all invalid chars (not from alphabet) with pad
func (c *cleanerImpl) clean(word string) string {
	word = c.normalize(word)
	return c.reg.ReplaceAllString(word, c.pad)
}

// wrap wraps word with wrapper string
func (c *cleanerImpl) wrap(word string) string {
	return c.wrapper + word + c.wrapper
}
