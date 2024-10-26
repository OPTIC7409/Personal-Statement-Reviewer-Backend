package JWT

import (
	"psr/database/queries"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(email string) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		Issuer:    "personal-statement-reviewer",
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("58 26 c4 c9 07 df ac 4f 19 40 a4 f9 e9 cc f0 b9  b6 0a ac 49 6d 29 14 13 cd 9b ae 6d e8 a2 2b 61  3c ac de 05 30 a5 1d 04 2a c4 d6 8b 0d 53 fa 2d  09 6d d7 4d 24 25 7e 18 12 f2 a3 0a 5d ac e9 de  ca de b5 bf 63 01 c4 91 a8 56 14 ce 38 4d 19 2e  0f 7a e6 b2 32 8d 63 29 a9 6e f3 3b 24 ac d0 22  b9 5e 69 22 5f c6 91 53 da 00 52 be 0d 5c 29 cb  92 d5 03 5f 8e 0f a6 4b 60 e1 da 66 77 dc 0a e1  b8 57 a0 e6 28 da 0d 01 22 84 71 16 26 bb 99 41  22 99 71 0d 73 2f c6 e7 7e 0c a2 99 d8 8b ba 66  e4 de 11 82 c8 8c c0 95 90 e9 89 17 ce bf e2 15  3e a1 f6 38 57 d5 68 d3 9b ed 12 cd ed a1 a8 fc  46 d8 42 3d 2b 35 23 03 11 19 dd 97 f0 c4 3a 49  8a 5c e9 a0 82 8b ce fd d5 54 ee 4a 85 00 28 3b  cd 10 8e cd c8 82 e6 1b 34 f4 0a 89 88 bf 63 69  6a c6 02 fb 17 27 6a ec 79 1e 34 2d e3 2a b2 bd"))

}

func ValidateJWT(token string) (string, error) {
	userID, err := queries.GetUserIDFromToken(token)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(userID), nil
}
