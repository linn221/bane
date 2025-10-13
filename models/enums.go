package models

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

type WordType string

const (
	WordTypeFuzz   WordType = "fuzz"
	WordTypeAttack WordType = "attack"
	WordTypeRegex  WordType = "regex"
)

func (t WordType) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(t))))
}

func (t *WordType) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("word type must be string")
	}
	switch str {
	case "fuzz":
		*t = WordTypeFuzz
	case "attack":
		*t = WordTypeAttack
	case "regex":
		*t = WordTypeRegex
	default:
		return errors.New("invalid word type")
	}
	return nil
}

type MyTime struct {
	time.Time
}

// MarshalGQL implements the graphql.Marshaler interface.
func (u MyTime) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(u.Format(time.DateOnly)))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (u *MyTime) UnmarshalGQL(v interface{}) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf("Uint must be a string")
	}

	return nil
}
