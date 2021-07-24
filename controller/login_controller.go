package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang-youtube-api/model"
	"golang-youtube-api/repository"
	"golang-youtube-api/security"
	"golang-youtube-api/service"
	"net/http"
	"os"
)

type loginController struct {
	userService service.UserService
	redis       security.Interface
	auth        security.TokenInterface
}

type LoginController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Refresh(c *gin.Context)
}

func NewLoginController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) LoginController {
	newLoginService := service.NewUserService(repo.User)
	return &loginController{newLoginService, redis, auth}
}

func (l *loginController) Login(c *gin.Context) {
	var user *model.User
	var tokenErr = map[string]string{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid json provided"})
		return
	}
	//validate request:
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateUser)
		return
	}
	u, userErr := l.userService.FindUserByEmailAndPassword(user)
	if userErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": userErr.Error()})
		return
	}
	ts, tErr := l.auth.CreateToken(u.UUID)
	if tErr != nil {
		tokenErr["token_error"] = tErr.Error()
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": tErr.Error()})
		return
	}
	saveErr := l.redis.CreateAuth(u.UUID, ts)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": saveErr.Error()})
		return
	}
	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	//userData["id"] = u.UUID
	//userData["first_name"] = u.FirstName
	//userData["last_name"] = u.LastName
	c.JSON(http.StatusOK, userData)
}

func (l *loginController) Logout(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := l.auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	//if the access token exist and it is still valid, then delete both the access token and the refresh token
	deleteErr := l.redis.DeleteTokens(metadata)
	if deleteErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": deleteErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (l *loginController) Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err})
		return
	}
	refreshToken := mapToken["refresh_token"]
	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//any error may be due to token expiration
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, "Cannot get uuid")
			return
		}
		userId := fmt.Sprint(claims["user_id"])
		userUUID, err := uuid.Parse(userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		//Delete the previous Refresh Token
		delErr := l.redis.DeleteRefresh(refreshUuid)
		if delErr != nil { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := l.auth.CreateToken(userUUID)
		if createErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": createErr.Error()})
			return
		}
		//save the tokens metadata to redis
		saveErr := l.redis.CreateAuth(userUUID, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": saveErr.Error()})
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
	}
}
