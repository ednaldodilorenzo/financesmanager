package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id interface{}, key *string, expTime *time.Duration) (*string, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)

	claims := tokenByte.Claims.(jwt.MapClaims)

	now := time.Now().UTC()

	claims["sub"] = id
	claims["exp"] = now.Add(*expTime).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(*key))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func ExtractSubContent(token, key *string) (*string, error) {
	tokenByte, err := jwt.Parse(*token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(*key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return nil, errors.New("invalid token")

	}

	result := fmt.Sprint(claims["sub"])

	return &result, nil
}
