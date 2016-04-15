package base

import (
	"math"
	"strings"

	"github.com/andyxning/shortme/constant"
)

// Int2String converts an unsigned 64bit integer to a string.
func Int2String(seq uint64) (shortURL string) {
	charSeq := []rune{}

	if seq != 0 {
		for seq != 0 {
			mod := seq % constant.BaseStringLength
			div := seq / constant.BaseStringLength
			charSeq = append(charSeq, rune(constant.BaseString[mod]))
			seq = div
		}
	} else {
		charSeq = append(charSeq, rune(constant.BaseString[seq]))
	}

	tmpShortURL := string(charSeq)
	shortURL = reverse(tmpShortURL)
	return
}

// String2Int converts a short URL string to an unsigned 64bit integer.
func String2Int(shortURL string) (seq uint64) {
	shortURL = reverse(shortURL)
	for index, char := range shortURL {
		base := uint64(math.Pow(float64(constant.BaseStringLength), float64(index)))
		seq += uint64(strings.Index(constant.BaseString, string(char)))*base
	}
	return
}