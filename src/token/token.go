package token

import (
	"fmt"
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

// secret key

var HMACSECRET = []byte("MYSECRET")

type TokenInfo struct {
	UserID string `json:"userid"`
}

type TokenClaims struct {
	UserID string `json:"userid"`
	NBF    time.Time `json:"nbf"`
	EXP    time.Time `json:"exp"`
	jwt.RegisteredClaims
}

// create accessToken and refreshToken

func Create(info TokenInfo) (string, string, error) {

	// 10 minutes
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": info.UserID,
		"nbf": time.Now(),
		"exp": time.Now().Add(time.Minute * 10),
	})

	// one day
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": info.UserID,
		"nbf": time.Now(),
		"exp": time.Now().AddDate(0, 0, 1),
	})

	// sign with secret
	// and get encoded token

	var err error
	accessTokenString, err := accessToken.SignedString(HMACSECRET)
	if err != nil{
		return "", "", errors.New("Could not create access token")
	}

	refreshTokenString, err := refreshToken.SignedString(HMACSECRET)
	if err != nil{
		return "", "", errors.New("Could not create refresh token")
	}

	return accessTokenString, refreshTokenString, nil
}

// verify token signature and expiration date
// returns nil if valid, err otherwise

func Validate(tokenString string) (error, *TokenClaims) {

	fmt.Println(tokenString)

	token, err := jwt.ParseWithClaims(
		tokenString, &TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {

		// validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return secret ([]byte)
		return HMACSECRET, nil
	})

	if err != nil {
		fmt.Println(err)
		return err, nil
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		fmt.Println(claims.UserID)
		fmt.Println(claims.NBF)
		fmt.Println(claims.EXP)
		fmt.Println("Expired?", !time.Now().Before(claims.EXP))

		if (time.Now().Before(claims.EXP)){
			return nil, claims
		}
	}

	return errors.New("Could not create refresh token"), nil
}
