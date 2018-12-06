package compression

// Encoder represents entity for encoding given posting list to byte array
type Encoder interface {
	// Encode encodes given positing list in byte array
	Encode(list []uint32) []byte
}

// Decoder represents entity for decoding given byte array to posting list
type Decoder interface {
	// Decode decodes given byte array to posting list
	Decode(bytes []byte) []uint32
}
