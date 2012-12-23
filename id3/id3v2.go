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
	h.Size = parseSize(data[6:])

	return h, nil
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
	for hasFrame(lreader, tagLen) {
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
			tags.Album = readString(lreader, size)
		case "track":
			tags.Track = readString(lreader, size)
		case "artist":
			tags.Artist = readString(lreader, size)
		case "name":
			tags.Name = readString(lreader, size)
		case "year":
			tags.Year = readString(lreader, size)
		case "disc":
			tags.Disc = readString(lreader, size)
		case "genre":
			tags.Genre = readGenre(lreader, size)
		case "length":
			tags.Length = readString(lreader, size)
		}
	}
	return tags, nil
}
