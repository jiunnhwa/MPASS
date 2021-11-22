package brute

import (
	"fmt"
	"mpass/logger"
	"mpass/util/file"
	"strings"
)

var BadWords []string

func init() {
	// Add an element to the filter.
	//BadWords = []string{"shit", "dumb", "fart"}
	BadWords, err := file.ReadAllLines("./data/badwords.txt") //[]string{"shit", "dumb", "fart"}
	fmt.Println(BadWords)                                     //otherwise compile complained not used!
	if err != nil {
		logger.Log("ERROR", 1, err.Error())
		return
	}
}

//HasBadWord test for membership in bad word list.
func HasBadWord(word string) bool {
	if InStr(word, BadWords) {
		return true
	}
	return false
}

func InStr(str string, list []string) bool {
	ustr := strings.ToUpper(strings.TrimSpace(str))
	for _, v := range list {
		if strings.Contains(ustr, strings.ToUpper(strings.TrimSpace(v))) {
			return true
		}
	}
	return false
}
