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

package main

import (
	"fmt"
	"os"

	"github.com/bpowers/go-id3"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: %s [FILE]...\n", os.Args[0])
		return
	}

	for _, s := range os.Args[1:] {
		f, err := os.Open(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open %s: %s\n", s, err)
			return
		}
		defer f.Close()

		tags := id3.Read(f)
		if tags == nil {
			fmt.Fprintf(os.Stderr, "Could not read ID3 information from %s\n", s)
			return
		}

		fmt.Println(s)
		fmt.Printf("Header\t%#v\n", tags.Header)
		fmt.Printf("Name\t%s\n", tags.Name)
		fmt.Printf("Artist\t%s\n", tags.Artist)
		fmt.Printf("Album\t%s\n", tags.Album)
		fmt.Printf("Year\t%s\n", tags.Year)
		fmt.Printf("Track\t%s\n", tags.Track)
		fmt.Printf("Disc\t%s\n", tags.Disc)
		fmt.Printf("Genre\t%s\n", tags.Genre)
		fmt.Printf("Length\t%s\n", tags.Length)
		fmt.Println()
	}
}
