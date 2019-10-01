package jwt

import (
	"errors"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// authentication constant ...
const (
	ATExpiredTime     = 24  // hours = 1 days
	RTExpiredTime     = 168 // hours = 7 days
	ATCookieName      = "_AT_SOFAST_"
	RTCookieName      = "_RT_SOFAST_"
	AccessSigningKey  = "Acc35sT0k3N!@#)(*s0fa5T"
	RefreshSigningKey = "!@#)(*s0fa5TR3fr3sHT0k3N"
	BearerPrefix      = "Bearer "
)

// Payload type for JWT payload
type Payload struct {
	UserID    int
	UserRoles []UserRole
}

// UserRole struct for JWT claims
type UserRole struct {
	RoleID   int
	NPSN     string
	SchoolID int
}

// Claims ...
type Claims struct {
	Payload *Payload
	jwt.StandardClaims
}

// AuthToken ...
type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateAuthToken ...
func GenerateAuthToken(payload *Payload) (*AuthToken, error) {
	claims := Claims{Payload: payload}
	signingMethod := jwt.SigningMethodHS256

	accessToken, err := generateAccessToken(claims, signingMethod)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	refreshToken, err := generateRefreshToken(claims, signingMethod)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &AuthToken{accessToken, refreshToken}, nil
}

func generateAccessToken(claims Claims, signingMethod jwt.SigningMethod) (string, error) {
	claims.ExpiresAt = time.Now().Add(ATExpiredTime * time.Hour).Unix()
	token := jwt.NewWithClaims(signingMethod, claims)

	return token.SignedString([]byte(AccessSigningKey))
}

func generateRefreshToken(claims Claims, signingMethod jwt.SigningMethod) (string, error) {
	claims.ExpiresAt = time.Now().Add(RTExpiredTime * time.Hour).Unix()
	token := jwt.NewWithClaims(signingMethod, claims)

	return token.SignedString([]byte(RefreshSigningKey))
}

// ValidateBearer ...
func ValidateBearer(bearerToken string) *Payload {
	accessToken := parseBearer(bearerToken)
	if accessToken == "" {
		return nil
	}

	claims, err := authenticate(accessToken)
	if err == nil && claims != nil {
		return claims.Payload
	}

	return nil
}

func parseBearer(bearerToken string) string {
	if strings.HasPrefix(bearerToken, BearerPrefix) {
		return strings.TrimPrefix(bearerToken, BearerPrefix)
	}

	return bearerToken
}

func authenticate(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("invalid method")
			return nil, errors.New("invalid method")
		}
		return []byte(AccessSigningKey), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("invalid access token")
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

// RefreshToken ...
func RefreshToken(refreshToken string) (*AuthToken, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("invalid method")
			return nil, errors.New("invalid method")
		}
		return []byte(RefreshSigningKey), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("invalid refresh token")
		return nil, errors.New("invalid refresh token")
	}

	return GenerateAuthToken(claims.Payload)
}
