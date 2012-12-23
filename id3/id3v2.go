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
	"encoding/binary"
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

var ID3v22Tags = map[string]string{
	"TAL": "album",
	"TRK": "track",
	"TP1": "artist",
	"TT2": "name",
	"TYE": "year",
	"TPA": "disc",
	"TCO": "genre",
}

var ID3v23Tags = map[string]string{
	"TALB": "album",
	"TRCK": "track",
	"TPE1": "artist",
	"TIT2": "name",
	"TYER": "year",
	"TPOS": "disc",
	"TCON": "genre",
	"TLEN": "length",
}

var ID3v24Tags = map[string]string{
	"TALB": "album",
	"TRCK": "track",
	"TPE1": "artist",
	"TIT2": "name",
	"TDRC": "year",
	"TPOS": "disc",
	"TCON": "genre",
	"TLEN": "length",
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

// ID3 v2.2 uses 24-bit big endian frame sizes.
func parseID3v22FrameSize(reader *bufio.Reader) (int, error) {
	size, err := readBytes(reader, 3)
	if err != nil {
		return -1, err
	}
	return int(size[0])<<16 | int(size[1])<<8 | int(size[2]), nil
}

// ID3 v2.3 doesn't use sync-safe frame sizes: read in as a regular big endian number.
func parseID3v23FrameSize(reader *bufio.Reader) (int, error) {
	var size int32
	binary.Read(reader, binary.BigEndian, &size)
	return int(size), nil
}

// ID3 v2.4 uses sync-safe frame sizes similar to those found in the header.
func parseID3v24FrameSize(reader *bufio.Reader) (int, error) {
	size, err := readBytes(reader, 4)
	if err != nil {
		return -1, err
	}
	return int(parseSize(size)), nil
}

func parseID3v2File(reader *bufio.Reader) (*File, error) {
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

	file := new(File)
	file.Header = header
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
			file.Album = readString(lreader, size)
		case "track":
			file.Track = readString(lreader, size)
		case "artist":
			file.Artist = readString(lreader, size)
		case "name":
			file.Name = readString(lreader, size)
		case "year":
			file.Year = readString(lreader, size)
		case "disc":
			file.Disc = readString(lreader, size)
		case "genre":
			file.Genre = readGenre(lreader, size)
		case "length":
			file.Length = readString(lreader, size)
		}
	}
	return file, nil
}
