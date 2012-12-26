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

var ID3v22Tags = map[string]string{
	"TAL": "album",
	"TP1": "artist",
	"COM": "comments",
	"TCM": "composer",
	"TCR": "copyright",
	"TPA": "disc",
	"TCO": "genre",
	"TEN": "encodedby",
	"TSS": "encoder",
	"TLA": "language",
	"TMT": "media",
	"TOA": "originalartist",
	"TT2": "title",
	"TRK": "track",
	"TXT": "writer",
	"TYE": "year",
}

// ID3 v2.2 uses 24-bit big endian frame sizes.
func parseID3v22FrameSize(reader *bufio.Reader) (int, error) {
	size, err := readBytes(reader, 3)
	if err != nil {
		return -1, err
	}
	return int(size[0])<<16 | int(size[1])<<8 | int(size[2]), nil
}
