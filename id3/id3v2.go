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
	"strings"
	"unicode/utf16"
)

// A parsed ID3v2 header as defined in Section 3 of
// http://id3.org/id3v2.4.0-structure
type ID3v2Header struct {
	Version           int
	MinorVersion      int
	Unsynchronization bool
	Extended          bool
	Experimental      bool
	Footer            bool
	Size              int32
}

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
		if b&0x80 > 0 {
			fmt.Println("Size byte had non-zero first bit")
		}

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
func parseID3v2String(data []byte) string {
	var s string
	switch data[0] {
	case 0: // ISO-8859-1 text.
		s = ISO8859_1ToUTF8(data[1:])
		break
	case 1: // UTF-16 with BOM.
		s = string(utf16.Decode(toUTF16(data[1:])))
		break
	case 2: // UTF-16BE without BOM.
		panic("Unsupported text encoding UTF-16BE.")
	case 3: // UTF-8 text.
		s = string(data[1:])
		break
	default:
		// No encoding, assume ISO-8859-1 text.
		s = ISO8859_1ToUTF8(data)
	}
	return strings.TrimRight(s, "\u0000")
}

func readID3v2String(reader *bufio.Reader, c int) string {
	b, err := readBytes(reader, c)
	if err != nil {
		// FIXME: return an error
		return ""
	}
	return parseID3v2String(b)
}

func readID3v2Genre(reader *bufio.Reader, c int) string {
	b, err := readBytes(reader, c)
	if err != nil {
		// FIXME: return an error
		return ""
	}
	genre := parseID3v2String(b)
	return convertID3v1Genre(genre)
}

func parseID3v2File(reader *bufio.Reader) (*SimpleTags, error) {
	var parseSize func(*bufio.Reader) (int, error)
	var tagMap map[string]string
	var tagLen int

	// parse header and setup version specific functions/data
	header, err := parseID3v2Header(reader)
	if err != nil {
		return nil, fmt.Errorf("parseHeader: %s", err)
	}
	switch header.Version {
	case 2:
		parseSize = parseID3v22FrameSize
		tagMap = ID3v22Tags
		tagLen = 3
	case 3:
		parseSize = parseID3v23FrameSize
		tagMap = ID3v23Tags
		tagLen = 4
	case 4:
		parseSize = parseID3v24FrameSize
		tagMap = ID3v24Tags
		tagLen = 4
	default:
		return nil, fmt.Errorf("Unrecognized ID3v2 version: %d", header.Version)
	}

	tags := new(SimpleTags)
	lreader := bufio.NewReader(io.LimitReader(reader, int64(header.Size)))
	for hasID3v2Frame(lreader, tagLen) {
		b, err := readBytes(lreader, tagLen)
		if err != nil {
			return nil, fmt.Errorf("parseID3v2File: %s", err)
		}
		tag := string(b)
		size, err := parseSize(lreader)
		if err != nil {
			return nil, err
		}
		// skip frame flags (only present in 2.3 and v2.4)
		if header.Version == 3 || header.Version == 4 {
			skipBytes(lreader, 2)
		}
		id, ok := tagMap[tag]
		if ok != true {
			// skip over unknown tags
			skipBytes(lreader, size)
		}

		switch id {
		case "album":
			tags.Album = readID3v2String(lreader, size)
		case "track":
			tags.Track = readID3v2String(lreader, size)
		case "artist":
			tags.Artist = readID3v2String(lreader, size)
		case "name":
			tags.Name = readID3v2String(lreader, size)
		case "year":
			tags.Year = readID3v2String(lreader, size)
		case "disc":
			tags.Disc = readID3v2String(lreader, size)
		case "genre":
			tags.Genre = readID3v2Genre(lreader, size)
		case "length":
			tags.Length = readID3v2String(lreader, size)
		}
	}
	return tags, nil
}
