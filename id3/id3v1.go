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
)

func hasID3v1Tag(reader io.ReadSeeker) bool {
	origin, err := reader.Seek(0, 1)
	if err != nil {
		return false
	}
	_, err = reader.Seek(-128, 2)
	if err != nil {
		return false
	}
	buf := make([]byte, 3)
	num, err := reader.Read(buf)
	if err != nil || num != 3 {
		return false
	}
	reader.Seek(origin, 0)
	return string(buf) == "TAG"
}

func parseID3v1File(reader io.ReadSeeker) (*SimpleTags, error) {
	origin, err := reader.Seek(-128, 2)
	if err != nil {
		return nil, fmt.Errorf("seek failed")
	}
	buf := bufio.NewReader(reader)

	header, err := readBytes(buf, 3)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	if string(header) != "TAG" {
		return nil, fmt.Errorf("ID3v1 tag not found")
	}
	tags := new(SimpleTags)

	data, err := readBytes(buf, 30)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	tags.Name = strings.TrimRight(string(data), "\u0000")

	data, err = readBytes(buf, 30)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	tags.Artist = strings.TrimRight(string(data), "\u0000")

	data, err = readBytes(buf, 30)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	tags.Album = strings.TrimRight(string(data), "\u0000")

	data, err = readBytes(buf, 4)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	tags.Year = strings.TrimRight(string(data), "\u0000")

	data, err = readBytes(buf, 30)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	if data[28] == 0 {
		tags.Track = fmt.Sprint(data[29])
	}

	data, err = readBytes(buf, 1)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	if int(data[0]) > len(id3v1Genres) {
		tags.Genre = "Unspecified"
	} else {
		tags.Genre = id3v1Genres[int(data[0])]
	}
	reader.Seek(origin, 0)
	return tags, nil
}
