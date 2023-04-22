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
		if len(word) > 2 && len(word) < 20 {
			if _, is_stopword := stopwords[word]; !is_stopword {
				res = append(res, word)
			}
		}

	}
	return res
}

func Occurances(str string, doc string) int {
	strToks := Tokenize(str)
	docToks := Tokenize(doc)
	count := 0
	for _, tok := range strToks {
		for _, dTok := range docToks {
			if tok == dTok {
				count++
			}
		}
	}
	return count
}
