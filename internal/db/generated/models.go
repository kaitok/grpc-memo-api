// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package generated

import (
	"time"

	"github.com/google/uuid"
)

type Memo struct {
	ID        uuid.UUID
	UserID    string
	Content   string
	CreatedAt time.Time
}
