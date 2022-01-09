package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func GerarSenha(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerificarSenha(senhaEncriptada, senhaProvida string) bool {
	if bcrypt.CompareHashAndPassword([]byte(senhaEncriptada), []byte(senhaProvida)) == nil {
		return true
	}
	return false
}
