package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessTokenDuration  = 15
	RefreshTokenDuration = 24 * 31
	AccessTokenName      = "Access-Token"
	RefreshTokenName     = "Refresh-Token"
)

type Claims struct {
	ExpireAt time.Time
	UserID   string
	Username string
	IssuedAt time.Time
}

func (c *Claims) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"exp":      jwt.NewNumericDate(c.ExpireAt),
		"id":       c.UserID,
		"username": c.Username,
		"iss_at":   jwt.NewNumericDate(time.Now()),
	}
}

type Pair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newAccessTokenClaims(userID, username string) *Claims {
	claims := newClaims(userID, username)
	claims.ExpireAt = claims.ExpireAt.Add(AccessTokenDuration * time.Minute)
	return claims
}

func newRefreshTokenClaims(userID, username string) *Claims {
	claims := newClaims(userID, username)
	claims.ExpireAt = claims.ExpireAt.Add(RefreshTokenDuration * time.Hour)
	return claims
}

func newClaims(userID, username string) *Claims {
	return &Claims{
		ExpireAt: time.Now(),
		UserID:   userID,
		Username: username,
		IssuedAt: time.Now(),
	}
}
