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
type SimpleTags struct {
	Title   string
	Artist string
	Album  string
	Year   string
	Track  string
	Disc   string
	Genre  string
	Length string
}

// Parse stream for ID3 information. Returns nil if parsing failed or the
// input didn't contain ID3 information.
// NOTE: ID3v1 and appended ID3v2.x are not supported without the ability
// to seek in the input. Use ReadFile instead.
func Read(reader io.Reader) (*SimpleTags, error) {
	buf := bufio.NewReader(reader)
	if !hasID3v2Tag(buf) {
		return nil, fmt.Errorf("no id3 tags")
	}
	tags, err := parseID3v2File(buf)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// Parse seekable stream for ID3 information. Returns nil if ID3 tag is
// not found or parsing fails.
func ReadFile(reader io.ReadSeeker) (*SimpleTags, error) {
	buf := bufio.NewReader(reader)
	if hasID3v1Tag(reader) {
		tags, err := parseID3v1File(reader)
		if err != nil {
			return nil, err
		}
		return tags, err
	} else if hasID3v2Tag(buf) {
		tags, err := parseID3v2File(buf)
		if err != nil {
			return nil, fmt.Errorf("error parsing ID3v2 tags")
		}
		return tags, err
	}
	return nil, fmt.Errorf("no id3 tags")
}
