package controller

import (
	"encoding/json"
	env "example/api_go/internal/helpers/secrets"
	"time"

	// db "example/api_go/client"
	model "example/api_go/model"
	"fmt"

	"example/api_go/internal/helpers/pages"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	// "example/api_go/internal/logger"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/*
InitializeOAuthGoogle Function
*/
func (cont *Controller) InitializeOAuthGoogle() {
	oauthConfGl.ClientID = env.GoDotEnvVariable("GOOGLE_CLIENT_ID")
	oauthConfGl.ClientSecret = env.GoDotEnvVariable("GOOGLE_CLIENT_SECRET")
	oauthStateStringGl = ""
}

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     env.GoDotEnvVariable("GOOGLE_CLIENT_ID"),
		ClientSecret: env.GoDotEnvVariable("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:9090/api/v1/auth/callback-gl",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

/*
CallBackFromGoogle Function
*/
func (cont *Controller) CallBackFromGoogle(c *gin.Context) {

	state := c.Request.FormValue("state")
	if state != oauthStateStringGl {
		fmt.Println("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Request.FormValue("code")
	if code == "" {
		fmt.Println("Code not found..")
		c.Writer.Write([]byte("Code Not Found to provide AccessToken..\n"))
		//c.Writer([]byte("Code Not Found to provide AccessToken..\n"))
		reason := c.Request.FormValue("error_reason")
		if reason == "user_denied" {
			c.Writer.Write([]byte("User has denied Permission.."))
		}
	} else {
		token, err := oauthConfGl.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return
		}

		cont.api.TokenSetTime = time.Now()
		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    token.AccessToken,
			HttpOnly: true,
			Path:     "/",
			Secure:   true,
			Expires:  time.Now().Add(1200 * time.Second),
		}

		http.SetCookie(c.Writer, cookie)
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			fmt.Println("Get: " + err.Error() + "\n")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		defer resp.Body.Close()

		var googleUser model.IGoogleUSer
		json.NewDecoder(resp.Body).Decode(&googleUser)

		Itoken := model.AccessToken{
			Id:     googleUser.Id,
			Token:  token.AccessToken,
			Expire: token.Expiry,
			Email:  googleUser.Email,
		}

		u, err := cont.data.UserBy("email", googleUser.Email)
		if err != nil || u.Token == "" {
			fmt.Println("User not found by GoogleID:")
			user := model.User{
				Username:       googleUser.Email,
				Email:          googleUser.Email,
				CreatedAt:      time.Now(),
				Token:          token.AccessToken,
				AllowTokenAuth: true,
				DataAllowed:    true,
			}
			cont.data.InsertUser(user)
		}
		cont.data.DeleteAccessToken(googleUser.Email)
		cont.data.CreateAccessToken(&Itoken)
		cont.data.SetAccessTokenToUser(&Itoken)

		c.JSON(http.StatusOK, model.Msg{Text: "LoggedIn", Token: token.AccessToken})
		return
	}
}

func (cont *Controller) HandleMain(c *gin.Context) {
	c.Render(
		http.StatusOK, render.Data{
			ContentType: "text/html; charset=utf-8",
			Data:        []byte(pages.IndexPage),
		})
}

/*
HandleLogin Function
*/
func (cont *Controller) HandleGoogleLogin(c *gin.Context) {
	URL, err := url.Parse(oauthConfGl.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: " + err.Error())
	}

	parameters := url.Values{}
	parameters.Add("client_id", oauthConfGl.ClientID)
	parameters.Add("scope", strings.Join(oauthConfGl.Scopes, " "))
	parameters.Add("redirect_uri", oauthConfGl.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateStringGl)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	c.Redirect(http.StatusTemporaryRedirect, url)
}
