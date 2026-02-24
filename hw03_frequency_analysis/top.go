package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type keyValue struct {
	key   string
	value uint
}

func Top10(str string) []string {
	result := make([]string, 0, 10)
	resultMap := make(map[string]uint)
	words := strings.Fields(str)
	for _, word := range words {
		clearWord := clearWord(word)
		if clearWord == "" {
			continue
		}
		resultMap[clearWord]++
	}
	sliceMap := make([]keyValue, 0, 10)
	for k, v := range resultMap {
		sliceMap = append(sliceMap, keyValue{k, v})
	}
	sort.Slice(sliceMap, func(i, j int) bool {
		if sliceMap[i].value == sliceMap[j].value {
			return sliceMap[i].key < sliceMap[j].key
		}
		return sliceMap[i].value > sliceMap[j].value
	})
	for _, v := range sliceMap {
		if len(result) == 10 {
			break
		}
		result = append(result, v.key)
	}
	return result
}

func clearWord(word string) string {
	if word == "-" {
		return ""
	}
	word = strings.ToLower(word)
	re := regexp.MustCompile(`^[\p{P}]+|[\p{P}]+$`)
	return re.ReplaceAllString(word, "")
}
