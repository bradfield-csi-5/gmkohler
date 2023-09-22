package pkg

import (
	"bytes"
	"testing"
)

type test struct {
	decoded string
	encoded []byte
}

var tests = []test{
	{
		decoded: "8.8.8.8",
		encoded: []byte{
			0x01, 0x38,
			0x01, 0x38,
			0x01, 0x38,
			0x01, 0x38,
			0x00,
		},
	},
	{
		decoded: "127.0.0.1",
		encoded: []byte{
			0x03, 0x31, 0x32, 0x37,
			0x01, 0x30,
			0x01, 0x30,
			0x01, 0x31,
			0x00,
		},
	},
	{
		decoded: "bradfieldcs.com",
		encoded: []byte{
			0x0b, 0x62, 0x72, 0x61, 0x64, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x63, 0x73,
			0x03, 0x63, 0x6f, 0x6d,
			0x00,
		},
	},
}

func TestResourceName_Encode(t *testing.T) {
	for _, data := range tests {
		name := ResourceName(data.decoded)
		encoded, err := name.Encode()
		if err != nil {
			t.Fatalf("unexpected error in encoding: %v", err)
		}
		encodedLength := len(encoded)
		expectedLength := len(data.encoded)
		if encodedLength != expectedLength {
			t.Fatalf(
				"expected encoded length %d, got %d",
				expectedLength,
				encodedLength,
			)
		}
		for j, b := range encoded {
			if b != data.encoded[j] {
				t.Fatalf(
					"bytes do not match.  Expected %v, got %v",
					data.encoded,
					encoded,
				)
			}
		}
	}
}

func TestDecodeQuestionName(t *testing.T) {
	for _, data := range tests {
		decoded, err := decodeResourceName(bytes.NewReader(data.encoded))
		if err != nil {
			t.Fatalf("unexpected error decoding: %v", err)
		}
		if decoded != ResourceName(data.decoded) {
			t.Fatalf("expected decoded %s, got %s", data.decoded, data.encoded)
		}
	}
}
