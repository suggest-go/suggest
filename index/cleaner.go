package index

import (
	"regexp"
	"strings"
)

// Cleaner provides api for preparing word for index search
type Cleaner interface {
	// CleanAndWrap returns prepared "cleaned" string
	CleanAndWrap(word string) string
	// CleanAndLeftWrap returns left wrapped "cleaned" string
	CleanAndLeftWrap(word string) string
}

// cleanerImpl implements Cleaner interface
type cleanerImpl struct {
	reg  *regexp.Regexp
	pad  string
	wrap [2]string
}

// NewCleaner returns new Cleaner object
func NewCleaner(chars []rune, pad string, wrap [2]string) Cleaner {
	str := string(chars)
	reg, err := regexp.Compile("[^" + str + "]+")
	if err != nil {
		panic(err)
	}

	return &cleanerImpl{
		reg:  reg,
		pad:  pad,
		wrap: wrap,
	}
}

// Clean returns prepared "cleaned" string
func (c *cleanerImpl) CleanAndWrap(word string) string {
	return c.wrap[0] + c.clean(word) + c.wrap[1]
}

func (c *cleanerImpl) CleanAndLeftWrap(word string) string {
	return c.wrap[0] + c.clean(word)
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
// TODO fix me (problem with ReplaceAllString for $ replacement, $ - is reserved symbol)
func (c *cleanerImpl) clean(word string) string {
	word = c.normalize(word)
	return c.reg.ReplaceAllString(word, c.pad)
}
