package services

import (
	env "example/api_go/internal/helpers/secrets"
	// "example/api_go/internal/logger"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     env.GoDotEnvVariable("GOOGLE_CLIENT_ID"),
		ClientSecret: env.GoDotEnvVariable("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:9090/callback-gl",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	oauthStateStringGl = ""
)

/*
InitializeOAuthGoogle Function
*/
func InitializeOAuthGoogle() {
	oauthConfGl.ClientID = env.GoDotEnvVariable("GOOGLE_CLIENT_ID")
	oauthConfGl.ClientSecret = env.GoDotEnvVariable("GOOGLE_CLIENT_SECRET")
	oauthStateStringGl = "/internal/logger"
}

/*
HandleGoogleLogin Function
*/
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	HandleLogin(w, r, oauthConfGl, oauthStateStringGl)
}

/*
CallBackFromGoogle Function
*/
func CallBackFromGoogle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World!")

	state := r.FormValue("state")
	if state != oauthStateStringGl {
		fmt.Println("invalid oauth state, expected " + oauthStateStringGl + ", got " + state + "\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	// logger.Log.Info(code)

	if code == "" {
		fmt.Println("Code not found..")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		// User has denied access..
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		token, err := oauthConfGl.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Println("oauthConfGl.Exchange() failed with " + err.Error() + "\n")
			return
		}
		fmt.Println("TOKEN>> AccessToken>> " + token.AccessToken)
		fmt.Println("TOKEN>> Expiration Time>> " + token.Expiry.String())
		fmt.Println("TOKEN>> RefreshToken>> " + token.RefreshToken)

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			fmt.Println("Get: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ReadAll: " + err.Error() + "\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		fmt.Println("parseResponseBody: " + string(response) + "\n")

		w.Write([]byte("Hello, I'm protected\n"))
		w.Write([]byte(string(response)))
		return
	}
}
