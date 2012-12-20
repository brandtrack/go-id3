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

// Package id3 implements basic ID3 parsing for MP3 files.
//
// Instead of providing access to every single ID3 frame this package
// exposes only the ID3v2 header and a few basic fields such as the
// artist, album, year, etc...
package id3

import (
	"bufio"
	"fmt"
	"io"
)

// A parsed ID3 file with common fields exposed.
type File struct {
	Header *ID3v2Header

	Name   string
	Artist string
	Album  string
	Year   string
	Track  string
	Disc   string
	Genre  string
	Length string
}

// Parse the input for ID3 information. Returns nil if parsing failed or the
// input didn't contain ID3 information.
func Read(reader io.Reader) (*File, error) {
	buf := bufio.NewReader(reader)
	if !isID3Tag(buf) {
		return nil, fmt.Errorf("no id3 tags")
	}
	tags, err := parseID3v2File(buf)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func isID3Tag(reader *bufio.Reader) bool {
	data, err := reader.Peek(3)
	if len(data) < 3 || err != nil {
		return false
	}
	return data[0] == 'I' && data[1] == 'D' && data[2] == '3'
}

