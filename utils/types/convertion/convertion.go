package convertion

import (
	"database/sql"
	"fmt"
)

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func NullStringToBool(ns sql.NullString) bool {
	fmt.Println("str: " + ns.String)
	fmt.Println(ns.Valid)

	if ns.Valid {
		return ns.String == "true"
	}
	return false
}
