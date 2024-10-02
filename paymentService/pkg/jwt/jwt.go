package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"math"
	"net/http"
	"time"
)

type Tokenizer struct {
	secret []byte
}

func (t *Tokenizer) ParseToken(tokenS string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		return t.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}

}

func (t *Tokenizer) GeneratePair(userID, username string) (*Pair, error) {
	claims := newAccessTokenClaims(userID, username)
	accessToken, err := t.generateToken(claims)
	if err != nil {
		return nil, err
	}

	claims = newRefreshTokenClaims(userID, username)
	refreshToken, err := t.generateToken(claims)
	if err != nil {
		return nil, err
	}

	return &Pair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t *Tokenizer) generateToken(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.ToMapClaims())
	return token.SignedString(t.secret)
}

func (t *Tokenizer) ParseMapClaims(mapClaims jwt.MapClaims) *Claims {
	sec, dec := math.Modf(mapClaims["exp"].(float64))
	ExpireAt := time.Unix(int64(sec), int64(dec*(1e9)))

	sec, dec = math.Modf(mapClaims["iss_at"].(float64))
	IssuedAt := time.Unix(int64(sec), int64(dec*(1e9)))

	return &Claims{
		ExpireAt: ExpireAt,
		UserID:   mapClaims["id"].(string),
		Username: mapClaims["username"].(string),
		IssuedAt: IssuedAt,
	}
}

func (t *Tokenizer) PrepareCookies(pair *Pair) (*http.Cookie, *http.Cookie) {
	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = AccessTokenName
	accessTokenCookie.Path = "/"
	accessTokenCookie.Value = pair.AccessToken
	accessTokenCookie.Secure = true
	accessTokenCookie.HttpOnly = true

	refreshTokenCookie := new(http.Cookie)
	refreshTokenCookie.Name = RefreshTokenName
	refreshTokenCookie.Path = "/"
	refreshTokenCookie.Value = pair.RefreshToken
	refreshTokenCookie.Secure = true
	refreshTokenCookie.HttpOnly = true

	return accessTokenCookie, refreshTokenCookie
}

func NewTokenizer(secret string) Tokenizer {
	return Tokenizer{
		secret: []byte(secret),
	}
}
