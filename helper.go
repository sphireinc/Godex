package Godex

import (
	"bytes"
	"encoding/json"
)

func softDeleteWhereClause(whereClause string) string {
	clause := " deleted_at IS NULL"
	if len(whereClause) > 0 {
		clause = clause + " AND " + whereClause
	}
	return clause
}

func Pretty(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return ""
	}
	return prettyJSON.String()
}
