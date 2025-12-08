package services

import (
	"BLOG/models"
	"BLOG/utils"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
}

var TokenMap = make(map[string]string)

// 注册
func (s *UserService) Register(user *models.Users) (bool, string) {

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("err hash password", err)
		return false, "注册失败"
	}
	user.Password = string(hashedPassword)
	db := utils.DBUtil{}.Connect()

	var existUser models.Users
	db.Where("username = ?", user.Username).First(&existUser)
	if existUser.Id != 0 {
		return false, "用户已存在"

	}
	db.Where("email = ?", user.Email).First(&existUser)
	if existUser.Id != 0 {
		return false, "邮箱已存在"

	}
	if err := db.Create(user).Error; err != nil {
		fmt.Println("err create user", err)
		return false, "create user error"
	}
	return true, "success"

}

// 登录
func (s UserService) Login(username, password string) (bool, string) {
	var user models.Users
	db := utils.DBUtil{}.Connect()
	db.Where("username = ?", username).Find(&user)

	var storedUser models.Users
	if err := db.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		fmt.Println("err find user", err)
		return false, "用户不存在"
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(password)); err != nil {
		fmt.Println("err compare password", err)
		return false, "密码错误"
	}
	//生成 jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.Id,
		"username": storedUser.Username,
		"email":    storedUser.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("abcdefghijklmnopqr"))
	if err != nil {
		fmt.Println("err generate token", err)
		return false, "生成token失败"
	}
	TokenMap[storedUser.Username] = tokenString
	return true, tokenString

}

// paraseToken 解析token
func (s UserService) paraseToken(tokenString string) (bool, models.Users) {
	var user models.Users
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("abcdefghijklmnopq"), nil
	})
	if err != nil {
		fmt.Println("err parse token", err)
		return false, user
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := int(claims["id"].(float64))
		username := claims["username"].(string)
		email := claims["email"].(string)

		//判断token是否一样
		exitsToken := TokenMap[username]
		if tokenString != exitsToken {
			return false, user

		}
		return true, models.Users{Id: id, Username: username, Email: email}
	}
	return false, user

}

// 用户登出
func (s UserService) Logout(tokenString string) (bool,error {
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},error) {
		return []byte("abcdefghijklmnopq"), nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		// 判断token 是否一致
		existToken := TokenMap[username]
		if tokenString == existToken {
			delete(TokenMap, username)
			return ok, nil
		}
	}
	return false, nil
}
