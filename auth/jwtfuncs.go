package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))

func GerarToken(argsToClaim map[string]string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	for k, v := range argsToClaim {
		claims[k] = v
	}

	tokenString, err := token.SignedString(SECRET_KEY)

	if err != nil {
		return "", fmt.Errorf("Algo deu errado: %s\n", err.Error())
	}

	return tokenString, nil
}

func PegarToken(bearerToken string) (string, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", fmt.Errorf("Tipo de autorização diferente do esperado")
	}

	tokenString := strings.Split(bearerToken, " ")[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SECRET_KEY, nil
	})
	if err != nil {
		return "", fmt.Errorf("Token invalido: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id, ok := claims["_id"].(string)
		if !ok {
			return "", err
		}
		return id, nil
	}
	fmt.Println(claims)
	return "", fmt.Errorf("Token está expirado")
}
