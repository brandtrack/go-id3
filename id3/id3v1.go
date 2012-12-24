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
	"fmt"
	"io"
	"strings"
)

type ID3v1Frame struct {
	name   string
	length int
}

var ID3v1Frames = []ID3v1Frame{
	{"title", 30},
	{"artist", 30},
	{"album", 30},
	{"year", 4},
	{"comment", 30},
}

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

func readID3v1String(reader io.Reader, c int) (string, error) {
	data, err := readBytes(reader, c)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\u0000"), nil
}

func parseID3v1File(reader io.ReadSeeker) (map[string]string, error) {
	origin, err := reader.Seek(-128, 2)
	if err != nil {
		return nil, fmt.Errorf("seek failed")
	}

	// verify tag header
	header, err := readID3v1String(reader, 3)
	if err != nil || header != "TAG" {
		return nil, fmt.Errorf("could not parse ID3v1 tag")
	}

	tags := map[string]string{}

	// parse simple string frames
	for _, v := range ID3v1Frames {
		str, err := readID3v1String(reader, v.length)
		if err != nil {
			return nil, fmt.Errorf("read error")
		}
		tags[v.name] = str
	}

	// parse track number (if present)
	_, err = reader.Seek(-2, 1)
	if err != nil {
		return nil, fmt.Errorf("seek error")
	}
	data, err := readBytes(reader, 2)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	if data[0] == 0 {
		tags["track"] = fmt.Sprint(data[1])
	}

	// parse genre
	data, err = readBytes(reader, 1)
	if err != nil {
		return nil, fmt.Errorf("read error")
	}
	if int(data[0]) > len(id3v1Genres) {
		tags["genre"] = "Unspecified"
	} else {
		tags["genre"] = id3v1Genres[int(data[0])]
	}

	reader.Seek(origin, 0)
	return tags, nil
}
