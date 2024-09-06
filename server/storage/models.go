package storage

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"` // Hash the password for security
	CreatedAt time.Time `json:"created_at"`
}
