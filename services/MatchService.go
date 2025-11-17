package services

import (
	"regexp"

	"github.com/linn221/bane/graph/model"
)

type Matchable interface {
	Text() string
}

func MatchRegex(obj Matchable, regexString string) (*model.SearchResult, error) {

	re, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}
	// jsonString, err := json.Marshal(obj)
	// if err != nil {
	// 	return nil, err
	// }
	text := obj.Text()
	results := re.FindAllString(text, -1)
	count := len(results)
	result := model.SearchResult{
		Results: results,
		Count:   &count,
	}
	return &result, nil
}
