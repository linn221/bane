package services

import "regexp"

type Matchable interface {
	Text() string
}

func MatchRegex(text Matchable, regexString string) ([]string, error) {
	re, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}
	results := re.FindAllString(text.Text(), -1)
	return results, nil
}
