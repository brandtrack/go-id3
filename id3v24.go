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
)

// ID3 v2.4 uses sync-safe frame sizes similar to those found in the header.
func parseID3v24Size(reader *bufio.Reader) (int, error) {
	size, err := readBytes(reader, 4)
	if err != nil {
		return -1, err
	}
	return int(parseSize(size)), nil
}

func parseID3v24File(reader *bufio.Reader) (*File, error) {
	file := new(File)
	for hasFrame(reader, 4) {
		b, err := readBytes(reader, 4)
		if err != nil {
			return nil, err
		}
		id := string(b)
		size, err := parseID3v24Size(reader)
		if err != nil {
			return nil, err
		}

		// Skip over frame flags.
		skipBytes(reader, 2)

		switch id {
		case "TALB":
			file.Album = readString(reader, size)
		case "TRCK":
			file.Track = readString(reader, size)
		case "TPE1":
			file.Artist = readString(reader, size)
		case "TCON":
			file.Genre = readGenre(reader, size)
		case "TIT2":
			file.Name = readString(reader, size)
		case "TDRC":
			// TODO: implement timestamp parsing
			file.Year = readString(reader, size)
		case "TPOS":
			file.Disc = readString(reader, size)
		case "TLEN":
			file.Length = readString(reader, size)
		default:
			skipBytes(reader, size)
		}
	}
	return file, nil
}
