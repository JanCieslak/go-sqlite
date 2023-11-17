package varint

import (
	"encoding/binary"
	"io"
)

// TODO: Test

func Decode(r io.Reader) (int64, error) {
	var a0 byte
	err := binary.Read(r, binary.BigEndian, &a0)
	if err != nil {
		return 0, err
	}

	var nBytes int

	if a0 > 240 {
		panic("lul")
	}

	switch {
	case a0 <= 240:
		return int64(a0), nil
	case 241 <= a0 && a0 <= 248:
		var a1 byte
		err := binary.Read(r, binary.BigEndian, &a1)
		if err != nil {
			return 0, err
		}
		return 240 + 256*(int64(a0)-241) + int64(a1), nil
	case a0 == 249:
		var a1 byte
		err := binary.Read(r, binary.BigEndian, &a1)
		if err != nil {
			return 0, err
		}
		var a2 byte
		err = binary.Read(r, binary.BigEndian, &a1)
		if err != nil {
			return 0, err
		}
		return 2288 + 256*int64(a1) + int64(a2), nil
	case a0 == 250:
		nBytes = 3
	case a0 == 251:
		nBytes = 4
	case a0 == 252:
		nBytes = 5
	case a0 == 253:
		nBytes = 6
	case a0 == 254:
		nBytes = 7
	case a0 == 255:
		nBytes = 8
	}

	n := make([]byte, nBytes)
	_, err = r.Read(n)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(n)), nil
}

func Encode(n uint64) []byte {
	switch {
	case n <= 240:
		return []byte{byte(n)}
	case n <= 2287:
		//return []byte{byte(n)}
	case n <= 67823:
	case n <= 16777215:
	case n <= 4294967295:
	case n <= 1099511627775:
	case n <= 281474976710655:
	case n <= 72057594037927935:
	}
	panic("lul")
}
