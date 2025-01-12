package utilities

import "golang.org/x/crypto/bcrypt"

type PasswordUtils struct {
}

func NewPasswordUtils() *PasswordUtils {
	return &PasswordUtils{}
}

func (p *PasswordUtils) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (p *PasswordUtils) CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
