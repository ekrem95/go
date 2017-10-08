package main

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims type
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	token, err := jwt.ParseWithClaims(
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im15dXNlcm5hbWUiLCJleHAiOjE1MDc0Njk0NjgsImlzcyI6ImxvY2FsaG9zdDo5MDAwIn0.mbeVXNfT0045TkOMMUMhTbBYuAl7a_EElcvG32fqvv8",
		&Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method")
			}
			return []byte("secret"), nil
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Grab the tokens claims and pass it into the original request
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Println(claims)
	}
}

// sign token
// func main() {
// 	expireToken := time.Now().Add(time.Minute * 5).Unix()
//
// 	// We'll manually assign the claims but in production you'd insert values from a database
// 	claims := Claims{
// 		"myusername",
// 		jwt.StandardClaims{
// 			ExpiresAt: expireToken,
// 			Issuer:    "localhost:9000",
// 		},
// 	}
//
// 	// Create the token using your claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
// 	// Signs the token with a secret.
// 	signedToken, _ := token.SignedString([]byte("secret"))
//
// 	fmt.Println(signedToken)
// }
