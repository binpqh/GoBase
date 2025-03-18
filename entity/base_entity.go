package entity

import "time"

type Entity[TKey any, TUserId any] struct {
	ID        TKey      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeleteAt  time.Time `json:"delete_at"`
	UpdatedBy TUserId   `json:"updated_by"`
	DeletedBy TUserId   `json:"deleted_by"`
	CreatedBy TUserId   `json:"created_by"`
}
