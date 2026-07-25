package model

import (
	"time"

	"github.com/uptrace/bun"
)

type PostLike struct {
	bun.BaseModel `bun:"table:posts_likes,alias:pl"`

	ID           int       `bun:"id,pk,autoincrement"`
	PostID       int       `bun:"post_id,notnull"`
	IP           *string   `bun:"ip"`
	DeviceMemory *int      `bun:"device_memory"`
	CreatedAt    time.Time `bun:"created_at,notnull"`
	UpdatedAt    time.Time `bun:"updated_at,notnull"`
}
