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

type tokenClaims struct {
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
		"exp": time.Now().AddDate(0, 0, -1),
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
// returns true | false if valid | invalid

func Validate(tokenString string) bool {

	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// return secret ([]byte)
		return HMACSECRET, nil
	})

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		fmt.Println(claims.UserID)
		fmt.Println("Expired?", !time.Now().Before(claims.EXP))

		if (time.Now().Before(claims.EXP)){
			return true
		}
	} else {
		fmt.Println("err", err)
	}

	return false
}
