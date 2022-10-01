package controller

import (
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

var (
	keySign        *rsa.PrivateKey
	keyVerify      *rsa.PublicKey
	keyVerifyBytes []byte
)

// read attempts to read local file
func read(key string) []byte {
	data, err := os.ReadFile(key)
	if err != nil {
		panic(err)
	}
	return data
}

// initializes package wide variables
func init() {
	var keyPathPrivate = "./private.rsa"
	var keyPathPublic = "./public.rsa.pub"
	var err error

	fmt.Println("initializing jwt package - setting private/public keys")
	keySignBytes := read(keyPathPrivate)
	keySign, err = jwt.ParseRSAPrivateKeyFromPEM(keySignBytes)
	if err != nil {
		panic(err)
	}

	keyVerifyBytes = read(keyPathPublic)
	keyVerify, err = jwt.ParseRSAPublicKeyFromPEM(keyVerifyBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println("initializing jwt package - done")
}

// GetPublicKey returns public signing key
// @Summary     returns public signing key
// @Tags        JWT
// @Produce     json
// @Success     200 {object} PublicKey
// @Router      /jwt [get]
func (c *Controller) GetPublicKey(context *gin.Context) {
	fmt.Println("public key request")
	context.IndentedJSON(http.StatusOK, PublicKey{PublicKey: string(keyVerifyBytes[:])})
}

// ValidateToken verifies token signature and returns its claims
// @Summary      verifies token signature and returns its claims related attributes
// @Tags         JWT
// @Produce      json
// @Param        token path string true "JWT token"
// @Success      200   {object} User
// @Failure      400   {object} Message "unable to verify token signature"
// @Failure      404   {object} Message "unable to find token specified user"
// @Router       /jwt/{token} [get]
func (c *Controller) ValidateToken(context *gin.Context) {
	fmt.Println("retrieving input payload")
	tokenString := context.Param("token")

	fmt.Println("verifying token signature")
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return keyVerify, nil
	})
	if err != nil {
		fmt.Println("unable to verify token")
		context.IndentedJSON(http.StatusBadRequest, Message{Message: "invalid token"})
		return
	}

	fmt.Println("retrieving and returning user attributes", token.Claims)
	tokenClaims := token.Claims.(*tokenClaims)
	if userAttributes, ok := users[tokenClaims.Username]; ok {
		context.IndentedJSON(http.StatusOK, userAttributes)
		return
	}

	fmt.Println("unable to find records for given token")
	context.IndentedJSON(http.StatusNotFound, Message{Message: "no user found for specified token"})
}

// GenerateToken generates jwt token
// @Summary      generates jwt token and includes username in its claims
// @Tags         JWT
// @Produce      json
// @Param        payload body Payload true "Client's username"
// @Success      201   {object} JwtToken "successfully create user jwt"
// @Failure      400   {object} Message "unable to parse payload"
// @Failure      500   {object} Message "internal server error"
// @Router       /jwt [post]
func (c *Controller) GenerateToken(context *gin.Context) {
	var payload Payload
	if err := context.BindJSON(&payload); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid input payload"})
		return
	}

	fmt.Println("generating token")
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	fmt.Println("assigning token claims")
	token.Claims = &tokenClaims{Username: payload.Username}

	fmt.Printf("signing token %s\n", token.Claims)
	tokenSigned, err := token.SignedString(keySign)
	if err != nil {
		fmt.Printf("error while signing token%s\n", err)
		context.IndentedJSON(http.StatusInternalServerError, Message{Message: err.Error()})
		return
	}

	fmt.Println("token successfully signed - generating response")
	context.IndentedJSON(http.StatusCreated, gin.H{"token": tokenSigned})
}
