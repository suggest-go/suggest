package suggest

type Encoder interface {
	Encode(list PostingList) []byte
}

type Decoder interface {
	Decode(bytes []byte) PostingList
}

func VBEncoder() Encoder {
	return &vcb{}
}

func VBDecoder() Decoder {
	return &vcb{}
}

type vcb struct {}

func (v *vcb) Encode(list PostingList) []byte {
	bytes := make([]byte, 0)
	buf := make([]byte, 0, 8) // fixme

	for _, n := range list {
		buf = buf[:0]
		vbEncode(n, buf)
		bytes = append(bytes, buf...)
	}

	return bytes
}

func (v *vcb) Decode(bytes []byte) PostingList {
	return nil
}

func vbEncode(n int, bytes []byte) {
	for {
		bytes = append(bytes, byte(n & 0x7F))
		if n < 0x80 {
			break
		}

		n >>= 7
	}

	l := len(bytes)

	for i := l / 2 - 1; i >= 0; i-- {
		opp := l - 1 - i
		bytes[i], bytes[opp] = bytes[opp], bytes[i]
	}

	bytes = append(bytes, byte(n | 0x80))
}

