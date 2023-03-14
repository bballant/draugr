package words

import (
	"regexp"
	"strings"
)

var exists = struct{}{}
var stopwords = map[string]struct{}{}

func init() {
	_stopwords := [...]string{"it", "the", "to", "a",
		"and", "be", "that", "this", "is", "are", "so", "i"}
	for _, v := range _stopwords {
		stopwords[v] = exists
	}
}

func keyify(str string) string {
	str = strings.ToLower(str)
	r, _ := regexp.Compile("([a-z\\-]+)")
	return r.FindString(str)
}

func Tokenize(str string) []string {
	words := strings.Fields(str)
	res := []string{}
	for _, v := range words {
		word := keyify(v)
		if _, stopword := stopwords[word]; !stopword {
			res = append(res, word)
		}

	}
	return res
}
