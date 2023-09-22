package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

type IpAddr [4]byte

func IpAddrFromString(s string) (IpAddr, error) {
	bytesAsStrings := strings.Split(s, ".")
	bytes := [4]byte{}
	if len(bytesAsStrings) != 4 {
		return bytes, fmt.Errorf("could not parse ip %s", s)
	}

	for j, s := range bytesAsStrings {
		b, err := strconv.Atoi(s)
		if err != nil {
			return bytes, fmt.Errorf("error parsing %s to byte: %v", s, err)
		}
		if b&0xff != b {
			return bytes, fmt.Errorf("value %s is larger than one byte", s)
		}
		bytes[j] = byte(b)
	}

	return bytes, nil
}
