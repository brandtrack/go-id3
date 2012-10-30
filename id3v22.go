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
)

// ID3 v2.2 uses 24-bit big endian frame sizes.
func parseID3v22FrameSize(reader *bufio.Reader) (int, error) {
	size, err := readBytes(reader, 3)
	if err != nil {
		return -1, err
	}
	return int(size[0])<<16 | int(size[1])<<8 | int(size[2]), nil
}

func parseID3v22File(reader *bufio.Reader) (*File, error) {
	file := new(File)
	for hasFrame(reader, 3) {
		b, err := readBytes(reader, 3)
		if err != nil {
			return nil, fmt.Errorf("parseID3v22File: %s", err)
		}
		id := string(b)
		size, err := parseID3v22FrameSize(reader)
		if err != nil {
			return nil, err
		}

		switch id {
		case "TAL":
			file.Album = readString(reader, size)
		case "TRK":
			file.Track = readString(reader, size)
		case "TP1":
			file.Artist = readString(reader, size)
		case "TT2":
			file.Name = readString(reader, size)
		case "TYE":
			file.Year = readString(reader, size)
		case "TPA":
			file.Disc = readString(reader, size)
		case "TCO":
			file.Genre = readGenre(reader, size)
		default:
			skipBytes(reader, size)
		}
	}
	return file, nil
}
