package queries

import (
	"psr/database"
	statement "psr/types/personal_statement"
)

func GetUserStatements(userID int) ([]statement.PersonalStatement, error) {
	var statements []statement.PersonalStatement
	_, err := database.GetConnection().Exec(`
		SELECT * FROM personal_statements WHERE user_id = $1
	`, userID)
	return statements, err
}

func SaveStatement(statement statement.PersonalStatement) error {
	_, err := database.GetConnection().Exec(`
		INSERT INTO personal_statements (content)
		VALUES (?)
	`, statement.Content)
	return err
}
