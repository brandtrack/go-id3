ID3 Parsing For Go
==================

Andrew Scherkus
May 21, 2012


Introduction
------------

Simple ID3 parsing library for go based on the specs at www.id3.org.

It doesn't handle everything but at least gets the imporant bits like artist,
album, track, etc...


Usage
-----
Pass in a suitable io.ReadSeeker and away you go!

    package main
    import (
            "os"
            "fmt"
            "github.com/bobertlo/go-id3/id3"
    )
    
    func main() {
            f, err := os.Open("foo.mp3")
            if err != nil {
                    return
            }
            defer f.Close()
            tags, err := id3.ReadFile(f)
            if err != nil {
                    return
            }
            fmt.Println("Title: ", tags["title"])
            fmt.Println("Artist: ", tags["artist"])
    }


Examples
--------
An example tag reading program can be found under id3/tagreader.

    go get github.com/bobertlo/go-id3/tagreader
    tagreader path/to/file.mp3 [...]
