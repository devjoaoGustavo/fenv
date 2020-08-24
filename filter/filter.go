package filter

import (
	"strings"
)

func SplitValues(rawTerms *string) (result []string) {
	terms := strings.Split(*rawTerms, ",")
	terms = Compact(terms)
	for i := range terms {
		result = append(result, terms[i])
	}
	return
}

func Compact(input []string) (result []string) {
	result = input[:0]
	for _, item := range input {
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	return
}
