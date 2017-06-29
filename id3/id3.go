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
	"errors"
	"fmt"
	"io"
)

// ReadFile parses seekable stream for ID3 information. Returns nil if
// ID3 tag is not found or parsing fails.
func ReadFile(reader io.ReadSeeker) (map[string]string, error) {
	buf := bufio.NewReader(reader)

	tags, v2err := parseID3v2File(buf)
	v1Tags, v1err := parseID3v1File(reader)

	if v1err != nil && v2err != nil {
		return nil, errors.New("Could not parse ID3 tags")
	}

	// Merge both results, prioricing id3v2
	for k, v := range v1Tags {
		if _, ok := tags[k]; !ok {
			tags[k] = v
		}
	}

	if len(tags) == 0 {
		return nil, fmt.Errorf("no id3 tags")
	}

	return tags, nil
}
