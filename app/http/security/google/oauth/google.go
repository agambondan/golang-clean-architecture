package oauth

import (
	"github.com/gin-gonic/gin"
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
		c.Redirect(http.StatusBadRequest, "/")
		return
	}
	code := c.Request.FormValue("code")
	if code == "" {
		c.Writer.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := c.Request.FormValue("error_reason")
		if reason == "user_denied" {
			c.Writer.Write([]byte("User has denied Permission.."))
		}
		//User has denied access..
		c.Redirect(http.StatusUnauthorized, "/")
		return
	} else {
		tokFile := "./assets/docs/token.json"
		token := getTokenFromWeb(AuthConfig(), code)
		saveToken(tokFile, token)
		c.JSON(http.StatusOK, gin.H{"installed": token})
		return
	}
	return
}
