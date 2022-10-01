package controller

import (
	"fmt"
	"github.com/golang-jwt/jwt"
)

type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Title    string   `json:"title"`
	Books    []string `json:"books"`
}

var users = map[string]User{
	"doejohn": {
		Username: "doejohn",
		Email:    "john.doe@example.com",
		Name:     "John Bastard Doe",
		Title:    "Lord Commander",
	},
	"doemary": {
		Username: "doemary",
		Email:    "mary.doe@example.com",
		Name:     "Mary Marei Doe",
		Title:    "Lady Commander",
	},
}

type tokenClaims struct {
	*jwt.StandardClaims
	Username string `json:"username" binding:"required"`
}

func (c tokenClaims) Valid() error {
	//TODO implement me
	fmt.Println("revisit this logic")
	return nil
}

type Payload struct {
	Username string `json:"username" binding:"required"`
}

type PublicKey struct {
	PublicKey string `json:"publickey" biding:"required"`
}

type JwtToken struct {
	JWT string `json:"jwt" binding:"required"`
}
