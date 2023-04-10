package structure

type User struct {
	Id           int    `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}
