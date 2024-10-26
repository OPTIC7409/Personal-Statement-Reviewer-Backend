package queries

import (
	"psr/database"
	statement "psr/types/personal_statement"
)

func GetUserStatements(userID int) ([]statement.PersonalStatement, error) {
	var statements []statement.PersonalStatement
	rows, err := database.GetConnection().Query(`
		SELECT * FROM personal_statements WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s statement.PersonalStatement
		if err := rows.Scan(&s.ID, &s.UserID, &s.Content, &s.CreatedAt); err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}

	return statements, nil
}

func SaveStatement(userID int, statement statement.PersonalStatement) error {
	_, err := database.GetConnection().Exec(`
		INSERT INTO personal_statements (user_id, content)
		VALUES ($1, $2)
	`, userID, statement.Content)
	return err
}

func GetStatementByFID(id int) (statement.PersonalStatement, error) {
	var s statement.PersonalStatement
	err := database.GetConnection().QueryRow(`
		SELECT id, user_id, content, created_at FROM personal_statements WHERE id = $1
	`, id).Scan(&s.ID, &s.UserID, &s.Content, &s.CreatedAt)
	if err != nil {
		return s, err
	}
	return s, nil
}

func UpdateStatement(id int, statement statement.PersonalStatement) error {
	_, err := database.GetConnection().Exec(`UPDATE personal_statements SET content = $1 WHERE id = $2`, statement.Content, id)
	return err
}
