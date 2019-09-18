// Package analysis represents API to convert text into indexable/searchable tokens
package analysis

// Token is a string with an assigned and thus identified meaning
type Token = string

// Tokenizer performs splitting the given text on a sequence of tokens
type Tokenizer interface {
	// Splits the given text on a sequence of tokens
	Tokenize(text string) []Token
}

// TokenFilter is responsible for removing, modifiying and altering the given token flow
type TokenFilter interface {
	// Filter filters the given list with described behaviour
	Filter(list []Token) []Token
}
