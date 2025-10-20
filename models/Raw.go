package models

import "github.com/linn221/bane/mystructs"

// QueryResult carries raw SQL results in a generic form
type QueryResult struct {
	Rows  [][]string
	Count int
}

// ToMyStrings converts the internal rows to []mystring with sep and limit
func (qr *QueryResult) ToMyStrings(sep *string, limit *int) []mystructs.MyString {
	if qr == nil || len(qr.Rows) == 0 {
		return []mystructs.MyString{}
	}
	out := make([]mystructs.MyString, 0, len(qr.Rows))
	for _, row := range qr.Rows {
		// Join columns by sep (default ",")
		s := ","
		if sep != nil {
			s = *sep
		}
		content := ""
		for i, col := range row {
			if i > 0 {
				content += s
			}
			content += col
		}
		out = append(out, mystructs.MyString{Content: content})
	}
	// Apply limit after building list
	if limit != nil && *limit >= 0 {
		if *limit < len(out) {
			return out[:*limit]
		}
	}
	return out
}
