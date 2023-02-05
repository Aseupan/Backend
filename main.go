package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "831741109899-mru1nrm9rme3l5e7c9krv1g0sneqi686.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
	userData         []byte
)

func handleMain(c *gin.Context) {
	var htmlIndex = `<html>
<body>
	<a href="/login">Google Log In</a>
</body>
</html>`
	c.Writer.Write([]byte(htmlIndex))
}

func handleGoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func handleGoogleCallback(c *gin.Context) {
	content, err := getUserInfo(c.Query("state"), c.Query("code"))
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	userData = content
	c.Redirect(http.StatusTemporaryRedirect, "/status")
}

type googleProfile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"`
}

func handleStatus(c *gin.Context) {
	if userData == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "not logged in",
		})
		return
	}

	var user googleProfile
	err := json.Unmarshal(userData, &user)

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

func handleTestLoggedIn(c *gin.Context) {
	if userData == nil {
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

func handleLogout(c *gin.Context) {
	userData = nil
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
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

func main() {
	r := gin.Default()
	r.GET("/", handleMain)
	r.GET("/login", handleGoogleLogin)
	r.GET("/callback", handleGoogleCallback)
	r.GET("/status", handleStatus)
	r.GET("/TLI", handleTestLoggedIn)
	r.GET("/logout", handleLogout)

	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
