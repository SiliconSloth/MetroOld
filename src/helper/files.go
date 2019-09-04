package helper

import (
	"io/ioutil"
	"strings"
)

func GetFileContents(location string) string {
	file, err := ioutil.ReadFile(location)
	if err != nil {
		//fmt.Println("No ", lastWord(location, '/'), " file")
		return ""
	}
	return string(file)
}

func lastWord(sentence string, seperator rune) string {
	var index int

	for i, r := range sentence {
		if r == seperator {
			index = i
		}
	}

	var out strings.Builder

	for i := index + 1; i < len(sentence); i++ {
		out.WriteRune([]rune(sentence)[i])
	}

	return out.String()
}