package user

type User struct {
	ID           int
	Username     string
	PasswordHash string
}

type CreateUser struct {
	Username     string
	PasswordHash string
}
