package revision

import "time"

type Revision struct {
	ID          int       `json:"id"`
	StatementID int       `json:"statement_id"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}
