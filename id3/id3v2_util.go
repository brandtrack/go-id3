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
	"strconv"
	"strings"
	"unicode/utf16"
)

func hasID3v2Tag(reader *bufio.Reader) bool {
	data, err := reader.Peek(3)
	if err != nil || len(data) < 3 {
		return false
	}
	return string(data) == "ID3"
}

// Peeks at the buffer to see if there is a valid frame.
func hasID3v2Frame(reader *bufio.Reader, frameSize int) bool {
	data, err := reader.Peek(frameSize)
	if err != nil {
		return false
	}

	for _, c := range data {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func parseID3v2Header(reader *bufio.Reader) (*ID3v2Header, error) {
	h := new(ID3v2Header)
	data, err := readBytes(reader, 10)
	if err != nil {
		return nil, fmt.Errorf("parseHeader: %s", err)
	}

	h.Version = int(data[3])
	h.MinorVersion = int(data[4])
	h.Unsynchronization = data[5]&1<<7 != 0
	h.Extended = data[5]&1<<6 != 0
	h.Experimental = data[5]&1<<5 != 0
	h.Footer = data[5]&1<<4 != 0
	h.Size = parseID3v2Size(data[6:])

	return h, nil
}

// Sizes are stored big endian but with the first bit set to 0 and always ignored.
// Refer to section 3.1 of http://id3.org/id3v2.4.0-structure
func parseID3v2Size(data []byte) int32 {
	size := int32(0)
	for i, b := range data {
		//if b&0x80 > 0 {
		//	fmt.Println("Size byte had non-zero first bit")
		//}

		shift := uint32(len(data)-i-1) * 7
		size |= int32(b&0x7f) << shift
	}
	return size
}

// Parses a string from frame data. The first byte represents the encoding:
//   0x01  ISO-8859-1
//   0x02  UTF-16 w/ BOM
//   0x03  UTF-16BE w/o BOM
//   0x04  UTF-8
//
// Refer to section 4 of http://id3.org/id3v2.4.0-structure
func parseID3v2String(data []byte) (string, error) {
	var s string
	switch data[0] {
	case 0: // ISO-8859-1 text.
		s = ISO8859_1ToUTF8(data[1:])
		break
	case 1: // UTF-16 with BOM.
		utf, err := toUTF16(data[1:])
		if err != nil {
			return "", err
		}
		s = string(utf16.Decode(utf))
		break
	case 2: // UTF-16BE without BOM.
		return "", fmt.Errorf("Unsupported text encoding UTF-16BE.")
	case 3: // UTF-8 text.
		s = string(data[1:])
		break
	default:
		// No encoding, assume ISO-8859-1 text.
		s = ISO8859_1ToUTF8(data)
	}
	return strings.TrimRight(s, "\u0000"), nil
}

func readID3v2String(reader *bufio.Reader, c int) (string, error) {
	b, err := readBytes(reader, c)
	if err != nil {
		return "", err
	}
	return parseID3v2String(b)
}

// ID3v2.2 and ID3v2.3 use "(NN)" where as ID3v2.4 simply uses "NN" when
// referring to ID3v1 genres. The "(NN)" format is allowed to have trailing
// information.
//
// RX and CR are shorthand for Remix and Cover, respectively.
//
// Refer to the following documentation:
//   http://id3.org/id3v2-00          TCO frame
//   http://id3.org/id3v2.3.0         TCON frame
//   http://id3.org/id3v2.4.0-frames  TCON frame
func convertID3v1Genre(genre string) string {
	if genre == "RX" || strings.HasPrefix(genre, "(RX)") {
		return "Remix"
	}
	if genre == "CR" || strings.HasPrefix(genre, "(CR)") {
		return "Cover"
	}

	// Try to parse "NN" format.
	index, err := strconv.Atoi(genre)
	if err == nil {
		if index >= 0 && index < len(id3v1Genres) {
			return id3v1Genres[index]
		}
		return "Unknown"
	}

	// Try to parse "(NN)" format.
	index = 0
	_, err = fmt.Sscanf(genre, "(%d)", &index)
	if err == nil {
		if index >= 0 && index < len(id3v1Genres) {
			return id3v1Genres[index]
		}
		return "Unknown"
	}

	// Couldn't parse so it's likely not an ID3v1 genre.
	return genre
}

func readID3v2Genre(reader *bufio.Reader, c int) (string, error) {
	b, err := readBytes(reader, c)
	if err != nil {
		return "", err
	}
	genre, err := parseID3v2String(b)
	if err != nil {
		return "", err
	}
	return convertID3v1Genre(genre), nil
}
