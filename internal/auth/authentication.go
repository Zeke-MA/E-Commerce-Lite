package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("unable to hash password")
	}
	return string(hashedPass), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return fmt.Errorf("password does not match")
	}

	return nil
}

func GenerateRefreshToken() (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	refreshToken := hex.EncodeToString(randBytes)

	return refreshToken, err
}

func GenerateJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    "ecom-lite",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return signedToken, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	tokenParts := strings.SplitN(authHeader, " ", 2)

	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid token authorization header format")
	}

	return tokenParts[1], nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, createKeyFunc(tokenSecret))

	if err != nil {
		return uuid.Nil, err
	}

	expired, err := token.Claims.GetExpirationTime()

	if err != nil {
		return uuid.Nil, err
	}

	if time.Now().After(expired.Time) {
		return uuid.Nil, jwt.ErrTokenExpired
	}

	issuer, err := token.Claims.GetIssuer()

	if err != nil {
		return uuid.Nil, err
	}

	if issuer != "ecom-lite" {
		return uuid.Nil, jwt.ErrTokenInvalidIssuer
	}

	id, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.Nil, err
	}

	tokenID, err := uuid.Parse(id)

	if err != nil {
		return uuid.Nil, err
	}

	return tokenID, nil
}

func createKeyFunc(tokenSecret string) func(t *jwt.Token) (interface{}, error) {
	return func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}
}
