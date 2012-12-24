// Copyright 2011 Andrew Scherkus
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package id3

import (
	"bufio"
	"fmt"
	"io"
)

var skipBuffer []byte = make([]byte, 1024*4)

func ISO8859_1ToUTF8(data []byte) string {
	p := make([]rune, len(data))
	for i, b := range data {
		p[i] = rune(b)
	}
	return string(p)
}

func toUTF16(data []byte) []uint16 {
	if len(data) < 2 {
		panic("Sequence is too short too contain a UTF-16 BOM")
	}
	if len(data)%2 > 0 {
		// TODO: if this is UTF-16 BE then this is likely encoded wrong
		data = append(data, 0)
	}

	var shift0, shift1 uint
	if data[0] == 0xFF && data[1] == 0xFE {
		// UTF-16 LE
		shift0 = 0
		shift1 = 8
	} else if data[0] == 0xFE && data[1] == 0xFF {
		// UTF-16 BE
		shift0 = 8
		shift1 = 0
		panic("UTF-16 BE found!")
	} else {
		panic(fmt.Sprintf("Unrecognized UTF-16 BOM: 0x%02X%02X", data[0], data[1]))
	}

	s := make([]uint16, 0, len(data)/2)
	for i := 2; i < len(data); i += 2 {
		s = append(s, uint16(data[i])<<shift0|uint16(data[i+1])<<shift1)
	}
	return s
}

func readBytes(reader io.Reader, c int) ([]byte, error) {
	b := make([]byte, c)

	n, err := reader.Read(b)
	if err != nil {
		return nil, err
	}
	if n != c {
		return nil, fmt.Errorf("short read, %d/%d", n, c)
	}
	return b, nil
}

func skipBytes(reader *bufio.Reader, c int) {
	pos := 0
	for pos < c {
		end := c - pos
		if end > len(skipBuffer) {
			end = len(skipBuffer)
		}

		i, err := reader.Read(skipBuffer[0:end])
		pos += i
		if err != nil {
			panic(err)
		}
	}
}
