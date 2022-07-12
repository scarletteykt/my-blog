package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 11

func Hash(password string) (string, error) {
	var hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return password, err
	}
	return string(hashedPassword), nil
}

func Compare(hashedPwd string, plainPwd string) error {
	byteHash := []byte(hashedPwd)
	bytePassword := []byte(plainPwd)
	return bcrypt.CompareHashAndPassword(byteHash, bytePassword)
}
