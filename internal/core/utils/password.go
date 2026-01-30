package utils

import "golang.org/x/crypto/bcrypt"

type PasswordManager struct{}

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

func (pm *PasswordManager) GenerateSalt() (string, error) {
	salt, err := bcrypt.GenerateFromPassword([]byte(""), bcrypt.DefaultCost)
	return string(salt), err
}

func (pm *PasswordManager) HashPassword(password string, salt string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(hash), err
}

func (pm *PasswordManager) VerifyPassword(hash, password string, salt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
	return err == nil
}
