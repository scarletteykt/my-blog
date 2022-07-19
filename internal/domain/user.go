package domain

type User struct {
	ID           int    `json:"id,omitempty" db:"id"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"password" db:"password_hash"`
}
