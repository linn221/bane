package models

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

type UInt uint

// MarshalGQL implements the graphql.Marshaler interface.
func (u UInt) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(strconv.FormatUint(uint64(u), 10)))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
func (u *UInt) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("Uint must be a string")
	}
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid Uint: %w", err)
	}
	*u = UInt(val)
	return nil
}

type WordType string

const (
	WordTypeFuzz   WordType = "fuzz"
	WordTypeAttack WordType = "attack"
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
	default:
		return errors.New("invalid word type")
	}
	return nil
}
