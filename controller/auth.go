package controller

import (
	ml "example/api_go/model"
	"fmt"

	"example/api_go/internal/helpers/secrets"
	"example/api_go/model"
	"example/api_go/version"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func (cont *Controller) Login(c *gin.Context) {
	defer cont.authMutex.Unlock()
	cont.authMutex.Lock()

	// Get the JSON body and decode into credentials
	var creds model.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPError{Error: err.Error()})
		return
	}

	// TODO return only invalid credentials msg

	email := strings.TrimSpace(strings.ToLower(creds.Email))

	expectedUser, err := cont.data.UserBy("email", email)
	if err != nil {
		log.Error().Msgf("error getting user data: %s. Email: %s", err.Error(), email)
		c.JSON(http.StatusUnauthorized, model.HTTPError{
			Error: "Usuario o Password Incorrecto",
			Code:  model.HTTPErrorCodeInvalidUserOrPass})
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, then we can move ahead
	// if NOT, then we return an "Unauthorized" status
	//fmt.Println(expectedUser.UserID, expectedUser.Password, expectedUser.Email, creds.Password)
	if expectedUser.ID.Hex() == "" {
		log.Error().Msgf("error getting user data: UserID is empty. Email: %s", email)
		c.JSON(http.StatusUnauthorized, model.HTTPError{
			Error: "Usuario o Password Incorrecto",
			Code:  model.HTTPErrorCodeInvalidUserOrPass})
		return
	}

	if expectedUser.Password == "" {
		log.Error().Msgf("error getting user data: Password is empty. Email: %s", email)
		c.JSON(http.StatusUnauthorized, model.HTTPError{
			Error: "Usuario o Password Incorrecto",
			Code:  model.HTTPErrorCodeInvalidUserOrPass})
		return
	}

	// TODO remove check for plain password
	if expectedUser.Password != creds.Password {
		err = secrets.Check(creds.Password, expectedUser.Password)
		if err != nil {
			//c.JSON(http.StatusUnauthorized, model.HTTPError{Error: "Invalid credentials 2"})
			log.Error().Msgf("error checking password: %s. Email: %s", err.Error(), expectedUser.Email)
			c.JSON(http.StatusUnauthorized, model.HTTPError{
				Error: "Usuario o Password Incorrecto",
				Code:  model.HTTPErrorCodeInvalidUserOrPass})
			return
		}
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()

	cont.api.TokenSetTime = time.Now()

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 1200 seconds, the same as the cache
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		Expires:  time.Now().Add(1200 * time.Second),
	}
	if version.Environment() == "dev" { //nolint
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, model.Msg{Text: "LoggedIn"})
}

func (cont *Controller) SignupUser(c *gin.Context) { //nolint:funlen
	// setHeaders(c)

	var param ml.SignupUserModel
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPError{Error: err.Error()})
		return
	}

	if param.Email == "" {
		c.JSON(http.StatusBadRequest, model.HTTPError{Error: "error an Email should be provided"})
		return
	}

	if param.Password == "" {
		c.JSON(http.StatusBadRequest, model.HTTPError{Error: "Error a Password should be provided"})
		return
	}

	_, err := cont.data.UserBy("email", param.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, model.HTTPError{Error: "Email already in use"})
		return
	}

	__, err := cont.data.UserBy("username", param.UserName)
	if err == nil {
		c.JSON(http.StatusBadRequest, model.HTTPError{Error: "Username already in use" + __.Username})
		return
	}

	pwd, err := secrets.Encrypt(param.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.HTTPError{Error: err.Error()})
		return
	}

	mUser := model.User{
		Username:  param.UserName,
		Email:     param.Email,
		Password:  pwd,
		CreatedAt: time.Now(),
	}

	cont.data.InsertUser(mUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.HTTPError{Error: "Error updating new password in db: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, mUser)
}

// getUser get the user from the session key
func (cont *Controller) getUser(c *gin.Context) (user ml.User, err error) {
	// This is set on the AuthRequired middleware
	response, ok := c.Get(responseKey)
	if !ok {
		return ml.User{}, fmt.Errorf("there is no user associated with the given key")
	}

	user, err = cont.data.UserBy("email", string(response.([]uint8)))
	if err != nil {
		return ml.User{}, fmt.Errorf("unable to get user from email")
	}

	return user, nil
}
