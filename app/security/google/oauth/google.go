package oauth

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

/*
HandleGoogleLogin Function
*/
func HandleGoogleLogin(c *gin.Context) {
	HandleLogin(c, AuthConfig(), "state-token")
}

/*
CallBackFromGoogle Function
*/
func CallBackFromGoogle(c *gin.Context) {
	state := c.Request.FormValue("state")
	if state != "state-token" {
		log.Println("invalid oauth state, expected " + "state-token" + ", got " + state + "\n")
		c.Redirect(http.StatusBadRequest, "/")
		return
	}
	code := c.Request.FormValue("code")
	if code == "" {
		log.Println("Code not found..")
		c.Writer.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := c.Request.FormValue("error_reason")
		if reason == "user_denied" {
			c.Writer.Write([]byte("User has denied Permission.."))
		}
		//User has denied access..
		c.Redirect(http.StatusUnauthorized, "/")
		return
	} else {
		tokFile := "token.json"
		token := getTokenFromWeb(AuthConfig(), code)
		saveToken(tokFile, token)
		c.JSON(http.StatusOK, gin.H{"installed": token})
		return
	}
	return
}
