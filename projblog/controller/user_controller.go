package controller

import (
	"BLOG/models"
	"BLOG/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func UserControllerInit(r *gin.Engine) {
	v1 := r.Group("user")
	{
		v1.POST("/register", Register)
	}
}
func Register(c *gin.Context) {
	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, msg := services.UserService{}.Register(user)
	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	}

}

// Login
func Login(c *gin.Context) {
	var user models.Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	ok, token := services.UserService{}.Login(user.Username, user.Password)
	if ok {
		c.JSON(http.StatusOK, gin.H{"message2": "Login successful",
			"token": token})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": token})
	}

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		// 验证token
		ok, user := services.UserService{}.ParaseToken(token)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			c.Abort()
			return
		}
		// 将用户信息保存到上下文中
		c.JSON(200, gin.H{
			"userId":   user.Id,
			"username": user.Username,
			"email":    user.Email,
		})
		c.Set("userId", user.Id)
		c.Set("username", user.Username)
		c.Set("email", user.Email)
		c.Next()
	}
}

// 登出
func Logout(c *gin.Context) {
	token := c.GetString("token")
	ok, err := services.UserService{}.Logout(token)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})

	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	userId, _ := c.Get("userId")
	username, _ := c.Get("username")
	mail, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{
		"message":  "Get user info successful",
		"userId":   userId,
		"username": username,
		"email":    mail})
}
