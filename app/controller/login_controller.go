package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-blog-api/app/lib"
	"go-blog-api/app/model"
	"go-blog-api/app/repository"
	"go-blog-api/app/security"
	"go-blog-api/app/service"
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
	Verify(c *gin.Context)
	VerifyRole(c *gin.Context)
}

func NewLoginController(repo *repository.Repositories, redis security.Interface, auth security.TokenInterface) LoginController {
	newLoginService := service.NewUserService(repo.User)
	return &loginController{newLoginService, redis, auth}
}

func (l *loginController) Login(c *gin.Context) {
	var user *model.User
	var userAPI *model.UserAPI
	if err := c.ShouldBindJSON(&userAPI); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), userAPI))
		return
	}
	_ = lib.Merge(userAPI, &user)
	//validate request:
	validateUser := user.Validate("login")
	if len(validateUser) > 0 {
		c.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("fill your empty field", "field can't empty", validateUser))
		return
	}
	cipherEncrypt, err := lib.CipherEncrypt([]byte(*user.Password), []byte(os.Getenv("CIPHER_KEY")))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("failed encrypt data", err.Error(), nil))
		return
	}
	u, userErr := l.userService.FindUserByEmailOrUsername(user)
	if userErr != nil {
		c.JSON(http.StatusNotFound, model.BuildErrorResponse("user not found", userErr.Error(), u))
		return
	}
	cipherDecrypt, err := lib.CipherDecrypt(cipherEncrypt, []byte(os.Getenv("CIPHER_KEY")))
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("failed decrypt data", err.Error(), nil))
		return
	}
	if *user.Password != fmt.Sprintf("%s", cipherDecrypt) {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", "password is wrong", nil))
		return
	}
	ts, tErr := l.auth.CreateToken(*u.ID)
	if tErr != nil {
		c.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("can't create token", tErr.Error(), nil))
		return
	}
	saveErr := l.redis.CreateAuth(*u.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusInternalServerError, model.BuildErrorResponse("can't save token to redis", saveErr.Error(), nil))
		return
	}
	userData := make(map[string]interface{})
	userData["access_token"] = ts.AccessToken
	userData["refresh_token"] = ts.RefreshToken
	c.JSON(http.StatusOK, userData)
}

func (l *loginController) Logout(c *gin.Context) {
	//check is the user is authenticated first
	metadata, err := l.auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), nil))
		return
	}
	//if the access token exist, and it is still valid, then delete both the access token and the refresh token
	deleteErr := l.redis.DeleteTokens(metadata)
	if deleteErr != nil {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("can't delete token", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, model.BuildResponse(true, "successfully logout", nil))
}

func (l *loginController) Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("invalid json", err.Error(), nil))
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
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("can't read refresh token", err.Error(), nil))
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), nil))
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("unauthorized", "cannot get uuid", nil))
			return
		}
		userId := fmt.Sprint(claims["user_id"])
		userID, err := uuid.Parse(userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, model.BuildErrorResponse("failed conver uuid", "Error occurred", nil))
			return
		}
		//Delete the previous Refresh Token
		delErr := l.redis.DeleteRefresh(refreshUuid)
		if delErr != nil { //if any goes wrong
			c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("can't delete refresh token", delErr.Error(), nil))
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := l.auth.CreateToken(userID)
		if createErr != nil {
			c.JSON(http.StatusForbidden, model.BuildErrorResponse("can't create token", createErr.Error(), nil))
			return
		}
		//save the tokens' metadata to redis
		saveErr := l.redis.CreateAuth(userID, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, model.BuildErrorResponse("can't save token to redis", saveErr.Error(), nil))
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, model.BuildResponse(true, "success", tokens))
	} else {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", "refresh token expired", nil))
	}
}

func (l *loginController) Verify(c *gin.Context) {
	_, err := l.auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, model.BuildResponse(true, "success", nil))
}

func (l *loginController) VerifyRole(c *gin.Context) {
	accessDetails, err := l.auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("unauthorized", err.Error(), nil))
		return
	}
	findById, err := l.userService.FindById(&accessDetails.UserUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.BuildErrorResponse("user not found", err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, model.BuildResponse(true, "success", findById.Role.Name))
}
