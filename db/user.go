package db

import (
	"bcraft/api/structures"
)

const (
	createUserRecord        = "INSERT INTO users(username, password_hash) VALUES($1, $2) RETURNING id"
	getUserRecordByUsername = "SELECT id, username, password_hash FROM users WHERE username=$1"
	getUserRecordById       = "SELECT id, username, password_hash FROM users WHERE id=$1"
)

func GetUserRecordById(uid int) (*structure.User, error) {
	var user structure.User

	err := Conn.Get(&user, getUserRecordById, uid)
	return &user, err
}

func GetUserRecord(username string) (*structure.User, error) {
	var user structure.User

	err := Conn.Get(&user, getUserRecordByUsername, username)
	return &user, err
}

func CreateUserRecord(user structure.User) (int, error) {
	var id int
	row := Conn.QueryRow(createUserRecord, user.Username, user.PasswordHash)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
