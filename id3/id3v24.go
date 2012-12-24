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

var ID3v24Tags = map[string]string{
	"TALB": "album",
	"TRCK": "track",
	"TPE1": "artist",
	"TIT2": "title",
	"TDRC": "year",
	"TPOS": "disc",
	"TCON": "genre",
	"TLEN": "length",
}

// ID3 v2.4 uses sync-safe frame sizes similar to those found in the header.
func parseID3v24FrameSize(reader *bufio.Reader) (int, error) {
	size, err := readBytes(reader, 4)
	if err != nil {
		return -1, err
	}
	return int(parseID3v2Size(size)), nil
}
