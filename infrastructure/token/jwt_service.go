package token

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	GenerateToken(userID string, expiration time.Duration, role string) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}

type jwtService struct {
	secretKey []byte
	issuer    string
}


func NewJWTService() JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET no está configurado")
	}
	return &jwtService{
		secretKey: []byte(secret),
		issuer:    "auth-service",
	}
}

func (s *jwtService) GenerateToken(userID string, expiration time.Duration, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iss":  s.issuer,
		"exp":  time.Now().Add(expiration).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secretKey)
}

func (s *jwtService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	}
	return nil, errors.New("token inválido o expirado")
}
