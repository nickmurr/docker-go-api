package model

import (
	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	jwtSecret = []byte("secret-go-api")
)

// User ...
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
	)

}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	return nil
}

// Sanitize
// Removing hidden field from response
func (u *User) Sanitize() {
	u.Password = ""
}

// Compare Password
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// func (u *User) TokenBack(mySigningKey string) string {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := make(jwt.MapClaims)
//
// 	// Устанавливаем набор параметров для токена
// 	claims["user"] = true
// 	claims["name"] = u.Email
// 	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
// 	token.Claims = claims
//
// 	// Подписываем токен нашим секретным ключем
// 	tokenString, _ := token.SignedString("secret")
// 	return tokenString
// }
type Claims struct {
	ID    int    `json:"user_id"`
	Email string `json:"user_email"`
	jwt.StandardClaims
}

func (u *User) TokenBack(mySigningKey []byte) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID:    u.ID,
		Email: u.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func CheckJwtToken(token string) (int, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, err
		}
		return 0, err
	}
	if !tkn.Valid {
		return 0, err
	}

	return claims.ID, nil
}
