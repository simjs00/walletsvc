package usecase

import (
	"github.com/golang-jwt/jwt"
)
var SampleSecretKey = []byte("SecretYouShouldHide")
type Account struct{
	 userDatabase map[string] bool
}

func InitAccountUsecase(UserData map[string]bool) *Account{
	return &Account{
		userDatabase: UserData,
	}
}

func (u *Account) CreateUser(customerXid string){
	u.userDatabase[customerXid]=true
}

func(u *Account) GetUser(customerXid string) bool{
	return u.userDatabase[customerXid]
}

func(u *Account) GenerateToken(customerXid string) (string,error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = u.GetUser(customerXid)
	claims["user"] = customerXid
	tokenString, err := token.SignedString(SampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString,nil
}