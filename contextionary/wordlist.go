/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 - 2018 Weaviate. All rights reserved.
 * LICENSE: https://github.com/creativesoftwarefdn/weaviate/blob/develop/LICENSE.md
 * AUTHOR: Bob van Luijt (bob@kub.design)
 * See www.creativesoftwarefdn.org for details
 * Contact: @CreativeSofwFdn / bob@kub.design
 */

package contextionary

// //// #include <string.h>
// //import "C"

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"syscall"
)

// Wordlist the actual wordlist and its metadata
type Wordlist struct {
	vectorWidth   uint64
	numberOfWords uint64
	metadata      map[string]interface{}

	file         os.File
	startOfTable int
	mmap         []byte
}

// LoadWordlist loads the wordlist or returns an error
func LoadWordlist(path string) (*Wordlist, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open the wordlist at %s: %+v", path, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("can't stat the wordlist at %s: %+v", path, err)
	}

	mmap, err := syscall.Mmap(int(file.Fd()), 0, int(fileInfo.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("can't mmap the file %s: %+v", path, err)
	}

	nrWordsBytes := mmap[0:8]
	vectorWidthBytes := mmap[8:16]
	metadataLengthBytes := mmap[16:24]

	nrWords := binary.LittleEndian.Uint64(nrWordsBytes)
	vectorWidth := binary.LittleEndian.Uint64(vectorWidthBytes)
	metadataLength := binary.LittleEndian.Uint64(metadataLengthBytes)

	metadataBytes := mmap[24 : 24+metadataLength]
	var metadata map[string]interface{}

	json.Unmarshal(metadataBytes, &metadata)

	// Compute beginning of word list lookup table.
	startOfTable := 24 + int(metadataLength)
	offset := 4 - (startOfTable % 4)
	startOfTable += offset

	return &Wordlist{
		vectorWidth:   vectorWidth,
		numberOfWords: nrWords,
		metadata:      metadata,
		startOfTable:  startOfTable,
		mmap:          mmap,
	}, nil
}

// GetNumberOfWords returns the number of words available
func (w *Wordlist) GetNumberOfWords() ItemIndex {
	return ItemIndex(w.numberOfWords)
}

// FindIndexByWord returns the index of a word
func (w *Wordlist) FindIndexByWord(_needle string) ItemIndex {
	var needle = string([]byte(_needle))
	needle += "\x00"

	var bytesNeedle = []byte(needle)

	var low ItemIndex
	high := ItemIndex(w.numberOfWords)

	for low <= high {
		midpoint := (low + high) / 2
		wordPtr := w.getWordPtr(midpoint)[0:len(bytesNeedle)]

		var cmp = bytes.Compare(bytesNeedle, wordPtr)

		if cmp == 0 {
			return midpoint
		} else if cmp < 0 {
			high = midpoint - 1
		} else {
			low = midpoint + 1
		}
	}

	return -1
}

func (w *Wordlist) getWordPtr(index ItemIndex) []byte {
	entryAddr := ItemIndex(w.startOfTable) + index*8
	wordAddressBytes := w.mmap[entryAddr : entryAddr+8]
	wordAddress := binary.LittleEndian.Uint64(wordAddressBytes)
	return w.mmap[wordAddress:]
}

func (w *Wordlist) getWord(index ItemIndex) string {
	ptr := w.getWordPtr(index)
	for i := 0; i < len(ptr); i++ {
		if ptr[i] == '\x00' {
			return string(ptr[0:i])
		}
	}

	return ""
}
