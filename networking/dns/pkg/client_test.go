package pkg

import "testing"

func TestIpFromString(t *testing.T) {
	type test struct {
		input string
		want  [4]byte
	}
	tests := []test{
		{
			input: "8.8.8.8",
			want:  [4]byte{0x08, 0x08, 0x08, 0x08},
		},
		{
			input: "255.255.255.255",
			want:  [4]byte{0xff, 0xff, 0xff, 0xff},
		},
		{
			input: "127.0.0.1",
			want:  [4]byte{0x7f, 0x00, 0x00, 0x01},
		},
	}
	for _, data := range tests {
		ipAddr, err := IpAddrFromString(data.input)
		if err != nil {
			t.Fatalf("Unexpected conversion failure")
		}
		if ipAddr != data.want {
			t.Logf("Expected %08x, got %08x", data.want, ipAddr)
			t.Fail()
		}
	}
}
