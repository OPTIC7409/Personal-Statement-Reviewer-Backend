package sanitize

import "strings"

func SanitizeQuery(query string) string {
	dangerousPatterns := []string{
		";", "--", "/*", "*/", "@@", "@",
		"char", "nchar", "varchar", "nvarchar",
		"alter", "begin", "cast", "create", "cursor", "declare", "delete", "drop", "end", "exec", "execute",
		"fetch", "insert", "kill", "select", "sys", "sysobjects", "syscolumns",
		"table", "update",
	}

	safeQuery := strings.ToLower(query)

	for _, pattern := range dangerousPatterns {
		safeQuery = strings.ReplaceAll(safeQuery, " "+pattern+" ", " ")
		safeQuery = strings.ReplaceAll(safeQuery, pattern+" ", " ")
		safeQuery = strings.ReplaceAll(safeQuery, " "+pattern, " ")
	}

	return safeQuery
}
