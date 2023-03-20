package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config
var OauthStateString = "pseudo-random"
var UserData []byte

func Init() {
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func HandleGoogleLogin(c *gin.Context) {
	url := GoogleOauthConfig.AuthCodeURL(OauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {
	content, err := GetUserInfo(c.Query("state"), c.Query("code"))
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	UserData = content
	c.Redirect(http.StatusTemporaryRedirect, "/status")
}

func GetUserInfo(state string, code string) ([]byte, error) {
	if state != OauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := GoogleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed getting user info, status code: %d", response.StatusCode)
	}
	return contents, nil
}

func HandleMain(c *gin.Context) {
	var htmlIndex = `<html>
<body>
	<a href="/login">Google Log In</a>
</body>
</html>`
	c.Writer.Write([]byte(htmlIndex))
}

type GoogleProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"`
}

func HandleStatus(c *gin.Context) {
	if UserData == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "not logged in",
		})
		return
	}

	var user GoogleProfile
	err := json.Unmarshal(UserData, &user)

	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   user,
	})
}

func HandleTestLoggedIn(c *gin.Context) {
	if UserData == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "not logged in",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "you are logged in",
	})
}

func HandleLogout(c *gin.Context) {
	UserData = nil
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
