package bloom

import (
	"fmt"
	"mpass/logger"
	"mpass/util/file"
	"strings"

	"github.com/yourbasic/bloom"
)

var blocklist *bloom.Filter

func init() {
	// A Bloom filter with room for at most 100000 elements.
	// The error rate for the filter is less than 1/200.
	// Add an element to the filter.
	blocklist = bloom.New(10000, 200)
	// Add an element to the filter.
	words, err := file.ReadAllLines("./data/badwords.txt") //[]string{"shit", "dumb", "fart"}
	fmt.Println(words)
	if err != nil {
		logger.Log("ERROR", 1, err.Error())
		return
	}
	for _, v := range words {
		blocklist.Add(v)
	}

}

//HasBadWord test for membership in bad word list.
func HasBadWord(word string) bool {
	for _, v := range strings.Fields(strings.ToLower(strings.TrimSpace(word))) {
		if blocklist.Test(v) {
			return true
		}
	}
	return false
}
