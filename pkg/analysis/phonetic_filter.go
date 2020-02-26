package analysis

import "github.com/suggest-go/phonetic"

type phoneticFilter struct {
	encoder phonetic.Encoder
}

func NewPhoneticFilter(encoder phonetic.Encoder) TokenFilter {
	return &phoneticFilter{
		encoder: encoder,
	}
}

// Filter filters the given list with described behaviour
func (p *phoneticFilter) Filter(list []Token) []Token {
	count := 0

	for _, token := range list {
		hash, err := p.encoder.Encode(token)

		if err != nil {
			panic(err) // TODO handle it
		}

		if len(hash) > 0 {
			list[count] = hash
			count++
		}
	}

	return list[:count]
}
